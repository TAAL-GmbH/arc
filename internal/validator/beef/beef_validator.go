package beef

import (
	"context"
	"errors"
	"fmt"

	"github.com/bitcoin-sv/arc/internal/beef"
	"github.com/bitcoin-sv/arc/internal/validator"
	"github.com/bitcoin-sv/arc/pkg/api"
	"github.com/libsv/go-bc"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript/interpreter"
	"github.com/ordishs/go-bitcoin"
)

type BeefValidator struct {
	policy *bitcoin.Settings
}

func New(policy *bitcoin.Settings) *BeefValidator {
	return &BeefValidator{
		policy: policy,
	}
}

func (v *BeefValidator) ValidateTransaction(ctx context.Context, beefTx *beef.BEEF, feeValidation validator.FeeValidation, scriptValidation validator.ScriptValidation) (*bt.Tx, error) {
	feeQuote := api.FeesToBtFeeQuote(v.policy.MinMiningTxFee)

	for _, btx := range beefTx.Transactions {
		// verify only unmined transactions
		if btx.IsMined() {
			continue
		}

		tx := btx.Transaction

		if err := validator.CommonValidateTransaction(v.policy, tx); err != nil {
			return tx, err
		}

		if feeValidation == validator.StandardFeeValidation {
			if err := standardCheckFees(tx, beefTx, feeQuote); err != nil {
				return tx, err
			}
		}

		if scriptValidation == validator.StandardScriptValidation {
			if err := validateScripts(tx, beefTx.Transactions); err != nil {
				return tx, err
			}
		}
	}

	if err := ensureAncestorsArePresentInBump(beefTx.GetLatestTx(), beefTx); err != nil {
		return beefTx.GetLatestTx(), validator.NewError(err, api.ErrStatusMinedAncestorsNotFound)
	}

	// TODO: verify merkle roots

	return nil, nil
}

func standardCheckFees(tx *bt.Tx, beefTx *beef.BEEF, feeQuote *bt.FeeQuote) error {
	expectedFees, err := validator.CalculateMiningFeesRequired(tx.SizeWithTypes(), feeQuote)
	if err != nil {
		return validator.NewError(err, api.ErrStatusFees)
	}

	inputSatoshis, outputSatoshis, err := calculateInputsOutputsSatoshis(tx, beefTx.Transactions)
	if err != nil {
		return validator.NewError(err, api.ErrStatusFees)
	}

	actualFeePaid := inputSatoshis - outputSatoshis

	if inputSatoshis < outputSatoshis {
		// force an error without wrong negative values
		actualFeePaid = 0
	}

	if actualFeePaid < expectedFees {
		err := fmt.Errorf("transaction fee of %d sat is too low - minimum expected fee is %d sat", actualFeePaid, expectedFees)
		return validator.NewError(err, api.ErrStatusFees)
	}

	return nil
}

func calculateInputsOutputsSatoshis(tx *bt.Tx, inputTxs []*beef.TxData) (uint64, uint64, error) {
	inputSum := uint64(0)

	for _, input := range tx.Inputs {
		inputParentTx := findParentForInput(input, inputTxs)

		if inputParentTx == nil {
			return 0, 0, errors.New("invalid parent transactions, no matching trasactions for input")
		}

		inputSum += inputParentTx.Transaction.Outputs[input.PreviousTxOutIndex].Satoshis
	}

	outputSum := tx.TotalOutputSatoshis()

	return inputSum, outputSum, nil
}

func validateScripts(tx *bt.Tx, inputTxs []*beef.TxData) error {
	for i, input := range tx.Inputs {
		inputParentTx := findParentForInput(input, inputTxs)
		if inputParentTx == nil {
			return validator.NewError(errors.New("invalid parent transactions, no matching trasactions for input"), api.ErrStatusUnlockingScripts)
		}

		err := checkScripts(tx, inputParentTx.Transaction, i)
		if err != nil {
			return validator.NewError(errors.New("invalid script"), api.ErrStatusUnlockingScripts)
		}
	}

	return nil
}

// TODO move to common
func checkScripts(tx, prevTx *bt.Tx, inputIdx int) error {
	input := tx.InputIdx(inputIdx)
	prevOutput := prevTx.OutputIdx(int(input.PreviousTxOutIndex))

	err := interpreter.NewEngine().Execute(
		interpreter.WithTx(tx, inputIdx, prevOutput),
		interpreter.WithForkID(),
		interpreter.WithAfterGenesis(),
	)

	return err
}

func ensureAncestorsArePresentInBump(tx *bt.Tx, beefTx *beef.BEEF) error {
	minedAncestors := make([]*beef.TxData, 0)

	for _, input := range tx.Inputs {
		if err := findMinedAncestorsForInput(input, beefTx.Transactions, &minedAncestors); err != nil {
			return err
		}
	}

	for _, tx := range minedAncestors {
		if !existsInBumps(tx, beefTx.BUMPs) {
			return errors.New("invalid BUMP - input mined ancestor is not present in BUMPs")
		}
	}

	return nil
}

func findMinedAncestorsForInput(input *bt.Input, ancestors []*beef.TxData, minedAncestors *[]*beef.TxData) error {
	parent := findParentForInput(input, ancestors)
	if parent == nil {
		return fmt.Errorf("invalid BUMP - cannot find mined parent for input %s", input.String())
	}

	if parent.IsMined() {
		*minedAncestors = append(*minedAncestors, parent)
		return nil
	}

	for _, in := range parent.Transaction.Inputs {
		err := findMinedAncestorsForInput(in, ancestors, minedAncestors)
		if err != nil {
			return err
		}
	}

	return nil
}

func findParentForInput(input *bt.Input, parentTxs []*beef.TxData) *beef.TxData {
	parentID := input.PreviousTxIDStr()

	for _, ptx := range parentTxs {
		if ptx.GetTxID() == parentID {
			return ptx
		}
	}

	return nil
}

func existsInBumps(tx *beef.TxData, bumps []*bc.BUMP) bool {
	bumpIdx := int(*tx.BumpIndex)
	txID := tx.GetTxID()

	if len(bumps) > bumpIdx {
		leafs := bumps[bumpIdx].Path[0]

		for _, lf := range leafs {
			if txID == *lf.Hash {
				return true
			}
		}
	}

	return false
}