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
	"github.com/ordishs/go-bitcoin"
)

type BeefValidator struct {
	policy     *bitcoin.Settings
	mrVerifier validator.MerkleVerifierI
}

func New(policy *bitcoin.Settings, mrVerifier validator.MerkleVerifierI) *BeefValidator {
	return &BeefValidator{
		policy:     policy,
		mrVerifier: mrVerifier,
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

	if feeValidation == validator.CumulativeFeeValidation {
		if err := cumulativeCheckFees(beefTx, feeQuote); err != nil {
			return beefTx.GetLatestTx(), err
		}
	}

	if err := ensureAncestorsArePresentInBump(beefTx.GetLatestTx(), beefTx); err != nil {
		return beefTx.GetLatestTx(), err
	}

	if vErr := verifyMerkleRoots(ctx, v.mrVerifier, beefTx); vErr != nil {
		return beefTx.GetLatestTx(), vErr
	}

	return nil, nil
}

func standardCheckFees(tx *bt.Tx, beefTx *beef.BEEF, feeQuote *bt.FeeQuote) *validator.Error {
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

func cumulativeCheckFees(beefTx *beef.BEEF, feeQuote *bt.FeeQuote) *validator.Error {
	cumulativeExpFee := uint64(0)
	cumulativePaidFee := uint64(0)

	for _, bTx := range beefTx.Transactions {
		tx := bTx.Transaction

		expFees, err := validator.CalculateMiningFeesRequired(tx.SizeWithTypes(), feeQuote)
		if err != nil {
			return validator.NewError(err, api.ErrStatusCumulativeFees)
		}

		cumulativeExpFee += expFees

		totalInputSatoshis, totalOutputSatoshis, err := calculateInputsOutputsSatoshis(tx, beefTx.Transactions)
		if err != nil {
			return validator.NewError(err, api.ErrStatusCumulativeFees)
		}

		cumulativePaidFee += (totalInputSatoshis - totalOutputSatoshis)
	}

	if cumulativeExpFee > cumulativePaidFee {
		err := fmt.Errorf("cumulative transaction fee of %d sat is too low - minimum expected fee is %d sat", cumulativePaidFee, cumulativeExpFee)
		return validator.NewError(err, api.ErrStatusCumulativeFees)
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

func validateScripts(tx *bt.Tx, inputTxs []*beef.TxData) *validator.Error {
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

func verifyMerkleRoots(ctx context.Context, v validator.MerkleVerifierI, beefTx *beef.BEEF) *validator.Error {
	mr, err := beef.CalculateMerkleRootsFromBumps(beefTx.BUMPs)
	if err != nil {
		return validator.NewError(err, api.ErrStatusCalculatingMerkleRoots)
	}

	unverifiedBlocks, err := v.Verify(ctx, mr)
	if err != nil {
		return validator.NewError(err, api.ErrStatusValidatingMerkleRoots)
	}

	if len(unverifiedBlocks) > 0 {
		err := fmt.Errorf("unable to verify BUMPs with block heights: %v", unverifiedBlocks)
		return validator.NewError(err, api.ErrStatusValidatingMerkleRoots)
	}

	return nil
}

func checkScripts(tx, prevTx *bt.Tx, inputIdx int) error {
	input := tx.InputIdx(inputIdx)
	prevOutput := prevTx.OutputIdx(int(input.PreviousTxOutIndex))

	return validator.CheckScript(tx, inputIdx, prevOutput)
}

func ensureAncestorsArePresentInBump(tx *bt.Tx, beefTx *beef.BEEF) *validator.Error {
	minedAncestors := make([]*beef.TxData, 0)

	for _, input := range tx.Inputs {
		if err := findMinedAncestorsForInput(input, beefTx.Transactions, &minedAncestors); err != nil {
			return validator.NewError(err, api.ErrStatusMinedAncestorsNotFound)
		}
	}

	for _, tx := range minedAncestors {
		if !existsInBumps(tx, beefTx.BUMPs) {
			err := errors.New("invalid BUMP - input mined ancestor is not present in BUMPs")
			return validator.NewError(err, api.ErrStatusMinedAncestorsNotFound)
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
