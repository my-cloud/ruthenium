// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package humanity

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// PohMetaData contains all meta data concerning the Poh contract.
var PohMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIArbitrator\",\"name\":\"_arbitrator\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_arbitratorExtraData\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"_registrationMetaEvidence\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_clearingMetaEvidence\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_submissionBaseDeposit\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"_submissionDuration\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"_renewalPeriodDuration\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"_challengePeriodDuration\",\"type\":\"uint64\"},{\"internalType\":\"uint256[3]\",\"name\":\"_multipliers\",\"type\":\"uint256[3]\"},{\"internalType\":\"uint64\",\"name\":\"_requiredNumberOfVouches\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_requestID\",\"type\":\"uint256\"}],\"name\":\"AddSubmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"_challengeID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumProofOfHumanity.Party\",\"name\":\"_party\",\"type\":\"uint8\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_contributor\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"AppealContribution\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIArbitrator\",\"name\":\"_arbitrator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_governor\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_submissionBaseDeposit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_submissionDuration\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_challengePeriodDuration\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_requiredNumberOfVouches\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_sharedStakeMultiplier\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_winnerStakeMultiplier\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_loserStakeMultiplier\",\"type\":\"uint256\"}],\"name\":\"ArbitratorComplete\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"_requestID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_challengeID\",\"type\":\"uint256\"}],\"name\":\"ChallengeResolved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIArbitrator\",\"name\":\"_arbitrator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"_disputeID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_metaEvidenceID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_evidenceGroupID\",\"type\":\"uint256\"}],\"name\":\"Dispute\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIArbitrator\",\"name\":\"_arbitrator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"_evidenceGroupID\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_party\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_evidence\",\"type\":\"string\"}],\"name\":\"Evidence\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"_challengeID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumProofOfHumanity.Party\",\"name\":\"_side\",\"type\":\"uint8\"}],\"name\":\"HasPaidAppealFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"_metaEvidenceID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_evidence\",\"type\":\"string\"}],\"name\":\"MetaEvidence\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_requestID\",\"type\":\"uint256\"}],\"name\":\"ReapplySubmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_requester\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_requestID\",\"type\":\"uint256\"}],\"name\":\"RemoveSubmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIArbitrator\",\"name\":\"_arbitrator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"_disputeID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_ruling\",\"type\":\"uint256\"}],\"name\":\"Ruling\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"_requestID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_challengeID\",\"type\":\"uint256\"}],\"name\":\"SubmissionChallenged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_voucher\",\"type\":\"address\"}],\"name\":\"VouchAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_voucher\",\"type\":\"address\"}],\"name\":\"VouchRemoved\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"string\",\"name\":\"_evidence\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"}],\"name\":\"addSubmission\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_submissionIDs\",\"type\":\"address[]\"},{\"internalType\":\"string[]\",\"name\":\"_evidence\",\"type\":\"string[]\"},{\"internalType\":\"string[]\",\"name\":\"_names\",\"type\":\"string[]\"}],\"name\":\"addSubmissionManually\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"}],\"name\":\"addVouch\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"arbitratorDataList\",\"outputs\":[{\"internalType\":\"contractIArbitrator\",\"name\":\"arbitrator\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"metaEvidenceUpdates\",\"type\":\"uint96\"},{\"internalType\":\"bytes\",\"name\":\"arbitratorExtraData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"arbitratorDisputeIDToDisputeData\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"challengeID\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"submissionID\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"challengePeriodDuration\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"},{\"internalType\":\"enumProofOfHumanity.Reason\",\"name\":\"_reason\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"_duplicateID\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_evidence\",\"type\":\"string\"}],\"name\":\"challengeRequest\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractIArbitrator\",\"name\":\"_arbitrator\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_arbitratorExtraData\",\"type\":\"bytes\"}],\"name\":\"changeArbitrator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"_submissionDuration\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"_renewalPeriodDuration\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"_challengePeriodDuration\",\"type\":\"uint64\"}],\"name\":\"changeDurations\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_governor\",\"type\":\"address\"}],\"name\":\"changeGovernor\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_loserStakeMultiplier\",\"type\":\"uint256\"}],\"name\":\"changeLoserStakeMultiplier\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"string\",\"name\":\"_registrationMetaEvidence\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_clearingMetaEvidence\",\"type\":\"string\"}],\"name\":\"changeMetaEvidence\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"_requiredNumberOfVouches\",\"type\":\"uint64\"}],\"name\":\"changeRequiredNumberOfVouches\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_sharedStakeMultiplier\",\"type\":\"uint256\"}],\"name\":\"changeSharedStakeMultiplier\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"_vouches\",\"type\":\"address[]\"},{\"internalType\":\"bytes[]\",\"name\":\"_signatures\",\"type\":\"bytes[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_expirationTimestamps\",\"type\":\"uint256[]\"}],\"name\":\"changeStateToPending\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_submissionBaseDeposit\",\"type\":\"uint256\"}],\"name\":\"changeSubmissionBaseDeposit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_winnerStakeMultiplier\",\"type\":\"uint256\"}],\"name\":\"changeWinnerStakeMultiplier\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_requestID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_duplicateID\",\"type\":\"address\"}],\"name\":\"checkRequestDuplicates\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"}],\"name\":\"executeRequest\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_challengeID\",\"type\":\"uint256\"},{\"internalType\":\"enumProofOfHumanity.Party\",\"name\":\"_side\",\"type\":\"uint8\"}],\"name\":\"fundAppeal\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"}],\"name\":\"fundSubmission\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getArbitratorDataListCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_requestID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_challengeID\",\"type\":\"uint256\"}],\"name\":\"getChallengeInfo\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"lastRoundID\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"challenger\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"disputeID\",\"type\":\"uint256\"},{\"internalType\":\"enumProofOfHumanity.Party\",\"name\":\"ruling\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"duplicateSubmissionIndex\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_requestID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_challengeID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_round\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_contributor\",\"type\":\"address\"}],\"name\":\"getContributions\",\"outputs\":[{\"internalType\":\"uint256[3]\",\"name\":\"contributions\",\"type\":\"uint256[3]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_requestID\",\"type\":\"uint256\"}],\"name\":\"getNumberOfVouches\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_requestID\",\"type\":\"uint256\"}],\"name\":\"getRequestInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"disputed\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"resolved\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"requesterLost\",\"type\":\"bool\"},{\"internalType\":\"enumProofOfHumanity.Reason\",\"name\":\"currentReason\",\"type\":\"uint8\"},{\"internalType\":\"uint16\",\"name\":\"nbParallelDisputes\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"lastChallengeID\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"arbitratorDataID\",\"type\":\"uint16\"},{\"internalType\":\"addresspayable\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"ultimateChallenger\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"usedReasons\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_requestID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_challengeID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_round\",\"type\":\"uint256\"}],\"name\":\"getRoundInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"appealed\",\"type\":\"bool\"},{\"internalType\":\"uint256[3]\",\"name\":\"paidFees\",\"type\":\"uint256[3]\"},{\"internalType\":\"enumProofOfHumanity.Party\",\"name\":\"sideFunded\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"feeRewards\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"}],\"name\":\"getSubmissionInfo\",\"outputs\":[{\"internalType\":\"enumProofOfHumanity.Status\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"submissionTime\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"index\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"registered\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"hasVouched\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"numberOfRequests\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"governor\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"}],\"name\":\"isRegistered\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"loserStakeMultiplier\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_requestID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_iterations\",\"type\":\"uint256\"}],\"name\":\"processVouches\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"string\",\"name\":\"_evidence\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"}],\"name\":\"reapplySubmission\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_evidence\",\"type\":\"string\"}],\"name\":\"removeSubmission\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"}],\"name\":\"removeSubmissionManually\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"}],\"name\":\"removeVouch\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"renewalPeriodDuration\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"requiredNumberOfVouches\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_disputeID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_ruling\",\"type\":\"uint256\"}],\"name\":\"rule\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"sharedStakeMultiplier\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"submissionBaseDeposit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"submissionCounter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"submissionDuration\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_evidence\",\"type\":\"string\"}],\"name\":\"submitEvidence\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"vouches\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"winnerStakeMultiplier\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_beneficiary\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_submissionID\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_requestID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_challengeID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_round\",\"type\":\"uint256\"}],\"name\":\"withdrawFeesAndRewards\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"withdrawSubmission\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// PohABI is the input ABI used to generate the binding from.
// Deprecated: Use PohMetaData.ABI instead.
var PohABI = PohMetaData.ABI

// Poh is an auto generated Go binding around an Ethereum contract.
type Poh struct {
	PohCaller     // Read-only binding to the contract
	PohTransactor // Write-only binding to the contract
	PohFilterer   // Log filterer for contract events
}

// PohCaller is an auto generated read-only Go binding around an Ethereum contract.
type PohCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PohTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PohTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PohFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PohFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PohSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PohSession struct {
	Contract     *Poh              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PohCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PohCallerSession struct {
	Contract *PohCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// PohTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PohTransactorSession struct {
	Contract     *PohTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PohRaw is an auto generated low-level Go binding around an Ethereum contract.
type PohRaw struct {
	Contract *Poh // Generic contract binding to access the raw methods on
}

// PohCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PohCallerRaw struct {
	Contract *PohCaller // Generic read-only contract binding to access the raw methods on
}

// PohTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PohTransactorRaw struct {
	Contract *PohTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPoh creates a new instance of Poh, bound to a specific deployed contract.
func NewPoh(address common.Address, backend bind.ContractBackend) (*Poh, error) {
	contract, err := bindPoh(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Poh{PohCaller: PohCaller{contract: contract}, PohTransactor: PohTransactor{contract: contract}, PohFilterer: PohFilterer{contract: contract}}, nil
}

// NewPohCaller creates a new read-only instance of Poh, bound to a specific deployed contract.
func NewPohCaller(address common.Address, caller bind.ContractCaller) (*PohCaller, error) {
	contract, err := bindPoh(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PohCaller{contract: contract}, nil
}

// NewPohTransactor creates a new write-only instance of Poh, bound to a specific deployed contract.
func NewPohTransactor(address common.Address, transactor bind.ContractTransactor) (*PohTransactor, error) {
	contract, err := bindPoh(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PohTransactor{contract: contract}, nil
}

// NewPohFilterer creates a new log filterer instance of Poh, bound to a specific deployed contract.
func NewPohFilterer(address common.Address, filterer bind.ContractFilterer) (*PohFilterer, error) {
	contract, err := bindPoh(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PohFilterer{contract: contract}, nil
}

// bindPoh binds a generic wrapper to an already deployed contract.
func bindPoh(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PohABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Poh *PohRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Poh.Contract.PohCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Poh *PohRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Poh.Contract.PohTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Poh *PohRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Poh.Contract.PohTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Poh *PohCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Poh.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Poh *PohTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Poh.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Poh *PohTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Poh.Contract.contract.Transact(opts, method, params...)
}

// ArbitratorDataList is a free data retrieval call binding the contract method 0xec0e71ba.
//
// Solidity: function arbitratorDataList(uint256 ) view returns(address arbitrator, uint96 metaEvidenceUpdates, bytes arbitratorExtraData)
func (_Poh *PohCaller) ArbitratorDataList(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Arbitrator          common.Address
	MetaEvidenceUpdates *big.Int
	ArbitratorExtraData []byte
}, error) {
	var out []interface{}
	err := _Poh.contract.Call(opts, &out, "arbitratorDataList", arg0)

	outstruct := new(struct {
		Arbitrator          common.Address
		MetaEvidenceUpdates *big.Int
		ArbitratorExtraData []byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Arbitrator = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.MetaEvidenceUpdates = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.ArbitratorExtraData = *abi.ConvertType(out[2], new([]byte)).(*[]byte)

	return *outstruct, err

}

// ArbitratorDataList is a free data retrieval call binding the contract method 0xec0e71ba.
//
// Solidity: function arbitratorDataList(uint256 ) view returns(address arbitrator, uint96 metaEvidenceUpdates, bytes arbitratorExtraData)
func (_Poh *PohSession) ArbitratorDataList(arg0 *big.Int) (struct {
	Arbitrator          common.Address
	MetaEvidenceUpdates *big.Int
	ArbitratorExtraData []byte
}, error) {
	return _Poh.Contract.ArbitratorDataList(&_Poh.CallOpts, arg0)
}

// ArbitratorDataList is a free data retrieval call binding the contract method 0xec0e71ba.
//
// Solidity: function arbitratorDataList(uint256 ) view returns(address arbitrator, uint96 metaEvidenceUpdates, bytes arbitratorExtraData)
func (_Poh *PohCallerSession) ArbitratorDataList(arg0 *big.Int) (struct {
	Arbitrator          common.Address
	MetaEvidenceUpdates *big.Int
	ArbitratorExtraData []byte
}, error) {
	return _Poh.Contract.ArbitratorDataList(&_Poh.CallOpts, arg0)
}

// ArbitratorDisputeIDToDisputeData is a free data retrieval call binding the contract method 0xdd254cd3.
//
// Solidity: function arbitratorDisputeIDToDisputeData(address , uint256 ) view returns(uint96 challengeID, address submissionID)
func (_Poh *PohCaller) ArbitratorDisputeIDToDisputeData(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (struct {
	ChallengeID  *big.Int
	SubmissionID common.Address
}, error) {
	var out []interface{}
	err := _Poh.contract.Call(opts, &out, "arbitratorDisputeIDToDisputeData", arg0, arg1)

	outstruct := new(struct {
		ChallengeID  *big.Int
		SubmissionID common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ChallengeID = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.SubmissionID = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// ArbitratorDisputeIDToDisputeData is a free data retrieval call binding the contract method 0xdd254cd3.
//
// Solidity: function arbitratorDisputeIDToDisputeData(address , uint256 ) view returns(uint96 challengeID, address submissionID)
func (_Poh *PohSession) ArbitratorDisputeIDToDisputeData(arg0 common.Address, arg1 *big.Int) (struct {
	ChallengeID  *big.Int
	SubmissionID common.Address
}, error) {
	return _Poh.Contract.ArbitratorDisputeIDToDisputeData(&_Poh.CallOpts, arg0, arg1)
}

// ArbitratorDisputeIDToDisputeData is a free data retrieval call binding the contract method 0xdd254cd3.
//
// Solidity: function arbitratorDisputeIDToDisputeData(address , uint256 ) view returns(uint96 challengeID, address submissionID)
func (_Poh *PohCallerSession) ArbitratorDisputeIDToDisputeData(arg0 common.Address, arg1 *big.Int) (struct {
	ChallengeID  *big.Int
	SubmissionID common.Address
}, error) {
	return _Poh.Contract.ArbitratorDisputeIDToDisputeData(&_Poh.CallOpts, arg0, arg1)
}

// ChallengePeriodDuration is a free data retrieval call binding the contract method 0x0082a36d.
//
// Solidity: function challengePeriodDuration() view returns(uint64)
func (_Poh *PohCaller) ChallengePeriodDuration(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Poh.contract.Call(opts, &out, "challengePeriodDuration")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// ChallengePeriodDuration is a free data retrieval call binding the contract method 0x0082a36d.
//
// Solidity: function challengePeriodDuration() view returns(uint64)
func (_Poh *PohSession) ChallengePeriodDuration() (uint64, error) {
	return _Poh.Contract.ChallengePeriodDuration(&_Poh.CallOpts)
}

// ChallengePeriodDuration is a free data retrieval call binding the contract method 0x0082a36d.
//
// Solidity: function challengePeriodDuration() view returns(uint64)
func (_Poh *PohCallerSession) ChallengePeriodDuration() (uint64, error) {
	return _Poh.Contract.ChallengePeriodDuration(&_Poh.CallOpts)
}

// CheckRequestDuplicates is a free data retrieval call binding the contract method 0x2e848506.
//
// Solidity: function checkRequestDuplicates(address _submissionID, uint256 _requestID, address _duplicateID) view returns(bool)
func (_Poh *PohCaller) CheckRequestDuplicates(opts *bind.CallOpts, _submissionID common.Address, _requestID *big.Int, _duplicateID common.Address) (bool, error) {
	var out []interface{}
	err := _Poh.contract.Call(opts, &out, "checkRequestDuplicates", _submissionID, _requestID, _duplicateID)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckRequestDuplicates is a free data retrieval call binding the contract method 0x2e848506.
//
// Solidity: function checkRequestDuplicates(address _submissionID, uint256 _requestID, address _duplicateID) view returns(bool)
func (_Poh *PohSession) CheckRequestDuplicates(_submissionID common.Address, _requestID *big.Int, _duplicateID common.Address) (bool, error) {
	return _Poh.Contract.CheckRequestDuplicates(&_Poh.CallOpts, _submissionID, _requestID, _duplicateID)
}

// CheckRequestDuplicates is a free data retrieval call binding the contract method 0x2e848506.
//
// Solidity: function checkRequestDuplicates(address _submissionID, uint256 _requestID, address _duplicateID) view returns(bool)
func (_Poh *PohCallerSession) CheckRequestDuplicates(_submissionID common.Address, _requestID *big.Int, _duplicateID common.Address) (bool, error) {
	return _Poh.Contract.CheckRequestDuplicates(&_Poh.CallOpts, _submissionID, _requestID, _duplicateID)
}

// GetArbitratorDataListCount is a free data retrieval call binding the contract method 0x90d7c13c.
//
// Solidity: function getArbitratorDataListCount() view returns(uint256)
func (_Poh *PohCaller) GetArbitratorDataListCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Poh.contract.Call(opts, &out, "getArbitratorDataListCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetArbitratorDataListCount is a free data retrieval call binding the contract method 0x90d7c13c.
//
// Solidity: function getArbitratorDataListCount() view returns(uint256)
func (_Poh *PohSession) GetArbitratorDataListCount() (*big.Int, error) {
	return _Poh.Contract.GetArbitratorDataListCount(&_Poh.CallOpts)
}

// GetArbitratorDataListCount is a free data retrieval call binding the contract method 0x90d7c13c.
//
// Solidity: function getArbitratorDataListCount() view returns(uint256)
func (_Poh *PohCallerSession) GetArbitratorDataListCount() (*big.Int, error) {
	return _Poh.Contract.GetArbitratorDataListCount(&_Poh.CallOpts)
}

// GetChallengeInfo is a free data retrieval call binding the contract method 0xd64240de.
//
// Solidity: function getChallengeInfo(address _submissionID, uint256 _requestID, uint256 _challengeID) view returns(uint16 lastRoundID, address challenger, uint256 disputeID, uint8 ruling, uint64 duplicateSubmissionIndex)
func (_Poh *PohCaller) GetChallengeInfo(opts *bind.CallOpts, _submissionID common.Address, _requestID *big.Int, _challengeID *big.Int) (struct {
	LastRoundID              uint16
	Challenger               common.Address
	DisputeID                *big.Int
	Ruling                   uint8
	DuplicateSubmissionIndex uint64
}, error) {
	var out []interface{}
	err := _Poh.contract.Call(opts, &out, "getChallengeInfo", _submissionID, _requestID, _challengeID)

	outstruct := new(struct {
		LastRoundID              uint16
		Challenger               common.Address
		DisputeID                *big.Int
		Ruling                   uint8
		DuplicateSubmissionIndex uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.LastRoundID = *abi.ConvertType(out[0], new(uint16)).(*uint16)
	outstruct.Challenger = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.DisputeID = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Ruling = *abi.ConvertType(out[3], new(uint8)).(*uint8)
	outstruct.DuplicateSubmissionIndex = *abi.ConvertType(out[4], new(uint64)).(*uint64)

	return *outstruct, err

}

// GetChallengeInfo is a free data retrieval call binding the contract method 0xd64240de.
//
// Solidity: function getChallengeInfo(address _submissionID, uint256 _requestID, uint256 _challengeID) view returns(uint16 lastRoundID, address challenger, uint256 disputeID, uint8 ruling, uint64 duplicateSubmissionIndex)
func (_Poh *PohSession) GetChallengeInfo(_submissionID common.Address, _requestID *big.Int, _challengeID *big.Int) (struct {
	LastRoundID              uint16
	Challenger               common.Address
	DisputeID                *big.Int
	Ruling                   uint8
	DuplicateSubmissionIndex uint64
}, error) {
	return _Poh.Contract.GetChallengeInfo(&_Poh.CallOpts, _submissionID, _requestID, _challengeID)
}

// GetChallengeInfo is a free data retrieval call binding the contract method 0xd64240de.
//
// Solidity: function getChallengeInfo(address _submissionID, uint256 _requestID, uint256 _challengeID) view returns(uint16 lastRoundID, address challenger, uint256 disputeID, uint8 ruling, uint64 duplicateSubmissionIndex)
func (_Poh *PohCallerSession) GetChallengeInfo(_submissionID common.Address, _requestID *big.Int, _challengeID *big.Int) (struct {
	LastRoundID              uint16
	Challenger               common.Address
	DisputeID                *big.Int
	Ruling                   uint8
	DuplicateSubmissionIndex uint64
}, error) {
	return _Poh.Contract.GetChallengeInfo(&_Poh.CallOpts, _submissionID, _requestID, _challengeID)
}

// GetContributions is a free data retrieval call binding the contract method 0x3a8363c2.
//
// Solidity: function getContributions(address _submissionID, uint256 _requestID, uint256 _challengeID, uint256 _round, address _contributor) view returns(uint256[3] contributions)
func (_Poh *PohCaller) GetContributions(opts *bind.CallOpts, _submissionID common.Address, _requestID *big.Int, _challengeID *big.Int, _round *big.Int, _contributor common.Address) ([3]*big.Int, error) {
	var out []interface{}
	err := _Poh.contract.Call(opts, &out, "getContributions", _submissionID, _requestID, _challengeID, _round, _contributor)

	if err != nil {
		return *new([3]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([3]*big.Int)).(*[3]*big.Int)

	return out0, err

}

// GetContributions is a free data retrieval call binding the contract method 0x3a8363c2.
//
// Solidity: function getContributions(address _submissionID, uint256 _requestID, uint256 _challengeID, uint256 _round, address _contributor) view returns(uint256[3] contributions)
func (_Poh *PohSession) GetContributions(_submissionID common.Address, _requestID *big.Int, _challengeID *big.Int, _round *big.Int, _contributor common.Address) ([3]*big.Int, error) {
	return _Poh.Contract.GetContributions(&_Poh.CallOpts, _submissionID, _requestID, _challengeID, _round, _contributor)
}

// GetContributions is a free data retrieval call binding the contract method 0x3a8363c2.
//
// Solidity: function getContributions(address _submissionID, uint256 _requestID, uint256 _challengeID, uint256 _round, address _contributor) view returns(uint256[3] contributions)
func (_Poh *PohCallerSession) GetContributions(_submissionID common.Address, _requestID *big.Int, _challengeID *big.Int, _round *big.Int, _contributor common.Address) ([3]*big.Int, error) {
	return _Poh.Contract.GetContributions(&_Poh.CallOpts, _submissionID, _requestID, _challengeID, _round, _contributor)
}

// GetNumberOfVouches is a free data retrieval call binding the contract method 0xdeb8f707.
//
// Solidity: function getNumberOfVouches(address _submissionID, uint256 _requestID) view returns(uint256)
func (_Poh *PohCaller) GetNumberOfVouches(opts *bind.CallOpts, _submissionID common.Address, _requestID *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Poh.contract.Call(opts, &out, "getNumberOfVouches", _submissionID, _requestID)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNumberOfVouches is a free data retrieval call binding the contract method 0xdeb8f707.
//
// Solidity: function getNumberOfVouches(address _submissionID, uint256 _requestID) view returns(uint256)
func (_Poh *PohSession) GetNumberOfVouches(_submissionID common.Address, _requestID *big.Int) (*big.Int, error) {
	return _Poh.Contract.GetNumberOfVouches(&_Poh.CallOpts, _submissionID, _requestID)
}

// GetNumberOfVouches is a free data retrieval call binding the contract method 0xdeb8f707.
//
// Solidity: function getNumberOfVouches(address _submissionID, uint256 _requestID) view returns(uint256)
func (_Poh *PohCallerSession) GetNumberOfVouches(_submissionID common.Address, _requestID *big.Int) (*big.Int, error) {
	return _Poh.Contract.GetNumberOfVouches(&_Poh.CallOpts, _submissionID, _requestID)
}

// GetRequestInfo is a free data retrieval call binding the contract method 0x6e112409.
//
// Solidity: function getRequestInfo(address _submissionID, uint256 _requestID) view returns(bool disputed, bool resolved, bool requesterLost, uint8 currentReason, uint16 nbParallelDisputes, uint16 lastChallengeID, uint16 arbitratorDataID, address requester, address ultimateChallenger, uint8 usedReasons)
func (_Poh *PohCaller) GetRequestInfo(opts *bind.CallOpts, _submissionID common.Address, _requestID *big.Int) (struct {
	Disputed           bool
	Resolved           bool
	RequesterLost      bool
	CurrentReason      uint8
	NbParallelDisputes uint16
	LastChallengeID    uint16
	ArbitratorDataID   uint16
	Requester          common.Address
	UltimateChallenger common.Address
	UsedReasons        uint8
}, error) {
	var out []interface{}
	err := _Poh.contract.Call(opts, &out, "getRequestInfo", _submissionID, _requestID)

	outstruct := new(struct {
		Disputed           bool
		Resolved           bool
		RequesterLost      bool
		CurrentReason      uint8
		NbParallelDisputes uint16
		LastChallengeID    uint16
		ArbitratorDataID   uint16
		Requester          common.Address
		UltimateChallenger common.Address
		UsedReasons        uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Disputed = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Resolved = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.RequesterLost = *abi.ConvertType(out[2], new(bool)).(*bool)
	outstruct.CurrentReason = *abi.ConvertType(out[3], new(uint8)).(*uint8)
	outstruct.NbParallelDisputes = *abi.ConvertType(out[4], new(uint16)).(*uint16)
	outstruct.LastChallengeID = *abi.ConvertType(out[5], new(uint16)).(*uint16)
	outstruct.ArbitratorDataID = *abi.ConvertType(out[6], new(uint16)).(*uint16)
	outstruct.Requester = *abi.ConvertType(out[7], new(common.Address)).(*common.Address)
	outstruct.UltimateChallenger = *abi.ConvertType(out[8], new(common.Address)).(*common.Address)
	outstruct.UsedReasons = *abi.ConvertType(out[9], new(uint8)).(*uint8)

	return *outstruct, err

}

// GetRequestInfo is a free data retrieval call binding the contract method 0x6e112409.
//
// Solidity: function getRequestInfo(address _submissionID, uint256 _requestID) view returns(bool disputed, bool resolved, bool requesterLost, uint8 currentReason, uint16 nbParallelDisputes, uint16 lastChallengeID, uint16 arbitratorDataID, address requester, address ultimateChallenger, uint8 usedReasons)
func (_Poh *PohSession) GetRequestInfo(_submissionID common.Address, _requestID *big.Int) (struct {
	Disputed           bool
	Resolved           bool
	RequesterLost      bool
	CurrentReason      uint8
	NbParallelDisputes uint16
	LastChallengeID    uint16
	ArbitratorDataID   uint16
	Requester          common.Address
	UltimateChallenger common.Address
	UsedReasons        uint8
}, error) {
	return _Poh.Contract.GetRequestInfo(&_Poh.CallOpts, _submissionID, _requestID)
}

// GetRequestInfo is a free data retrieval call binding the contract method 0x6e112409.
//
// Solidity: function getRequestInfo(address _submissionID, uint256 _requestID) view returns(bool disputed, bool resolved, bool requesterLost, uint8 currentReason, uint16 nbParallelDisputes, uint16 lastChallengeID, uint16 arbitratorDataID, address requester, address ultimateChallenger, uint8 usedReasons)
func (_Poh *PohCallerSession) GetRequestInfo(_submissionID common.Address, _requestID *big.Int) (struct {
	Disputed           bool
	Resolved           bool
	RequesterLost      bool
	CurrentReason      uint8
	NbParallelDisputes uint16
	LastChallengeID    uint16
	ArbitratorDataID   uint16
	Requester          common.Address
	UltimateChallenger common.Address
	UsedReasons        uint8
}, error) {
	return _Poh.Contract.GetRequestInfo(&_Poh.CallOpts, _submissionID, _requestID)
}

// GetRoundInfo is a free data retrieval call binding the contract method 0xa84dc70e.
//
// Solidity: function getRoundInfo(address _submissionID, uint256 _requestID, uint256 _challengeID, uint256 _round) view returns(bool appealed, uint256[3] paidFees, uint8 sideFunded, uint256 feeRewards)
func (_Poh *PohCaller) GetRoundInfo(opts *bind.CallOpts, _submissionID common.Address, _requestID *big.Int, _challengeID *big.Int, _round *big.Int) (struct {
	Appealed   bool
	PaidFees   [3]*big.Int
	SideFunded uint8
	FeeRewards *big.Int
}, error) {
	var out []interface{}
	err := _Poh.contract.Call(opts, &out, "getRoundInfo", _submissionID, _requestID, _challengeID, _round)

	outstruct := new(struct {
		Appealed   bool
		PaidFees   [3]*big.Int
		SideFunded uint8
		FeeRewards *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Appealed = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.PaidFees = *abi.ConvertType(out[1], new([3]*big.Int)).(*[3]*big.Int)
	outstruct.SideFunded = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	outstruct.FeeRewards = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetRoundInfo is a free data retrieval call binding the contract method 0xa84dc70e.
//
// Solidity: function getRoundInfo(address _submissionID, uint256 _requestID, uint256 _challengeID, uint256 _round) view returns(bool appealed, uint256[3] paidFees, uint8 sideFunded, uint256 feeRewards)
func (_Poh *PohSession) GetRoundInfo(_submissionID common.Address, _requestID *big.Int, _challengeID *big.Int, _round *big.Int) (struct {
	Appealed   bool
	PaidFees   [3]*big.Int
	SideFunded uint8
	FeeRewards *big.Int
}, error) {
	return _Poh.Contract.GetRoundInfo(&_Poh.CallOpts, _submissionID, _requestID, _challengeID, _round)
}

// GetRoundInfo is a free data retrieval call binding the contract method 0xa84dc70e.
//
// Solidity: function getRoundInfo(address _submissionID, uint256 _requestID, uint256 _challengeID, uint256 _round) view returns(bool appealed, uint256[3] paidFees, uint8 sideFunded, uint256 feeRewards)
func (_Poh *PohCallerSession) GetRoundInfo(_submissionID common.Address, _requestID *big.Int, _challengeID *big.Int, _round *big.Int) (struct {
	Appealed   bool
	PaidFees   [3]*big.Int
	SideFunded uint8
	FeeRewards *big.Int
}, error) {
	return _Poh.Contract.GetRoundInfo(&_Poh.CallOpts, _submissionID, _requestID, _challengeID, _round)
}

// GetSubmissionInfo is a free data retrieval call binding the contract method 0x97973043.
//
// Solidity: function getSubmissionInfo(address _submissionID) view returns(uint8 status, uint64 submissionTime, uint64 index, bool registered, bool hasVouched, uint256 numberOfRequests)
func (_Poh *PohCaller) GetSubmissionInfo(opts *bind.CallOpts, _submissionID common.Address) (struct {
	Status           uint8
	SubmissionTime   uint64
	Index            uint64
	Registered       bool
	HasVouched       bool
	NumberOfRequests *big.Int
}, error) {
	var out []interface{}
	err := _Poh.contract.Call(opts, &out, "getSubmissionInfo", _submissionID)

	outstruct := new(struct {
		Status           uint8
		SubmissionTime   uint64
		Index            uint64
		Registered       bool
		HasVouched       bool
		NumberOfRequests *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Status = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.SubmissionTime = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.Index = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.Registered = *abi.ConvertType(out[3], new(bool)).(*bool)
	outstruct.HasVouched = *abi.ConvertType(out[4], new(bool)).(*bool)
	outstruct.NumberOfRequests = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetSubmissionInfo is a free data retrieval call binding the contract method 0x97973043.
//
// Solidity: function getSubmissionInfo(address _submissionID) view returns(uint8 status, uint64 submissionTime, uint64 index, bool registered, bool hasVouched, uint256 numberOfRequests)
func (_Poh *PohSession) GetSubmissionInfo(_submissionID common.Address) (struct {
	Status           uint8
	SubmissionTime   uint64
	Index            uint64
	Registered       bool
	HasVouched       bool
	NumberOfRequests *big.Int
}, error) {
	return _Poh.Contract.GetSubmissionInfo(&_Poh.CallOpts, _submissionID)
}

// GetSubmissionInfo is a free data retrieval call binding the contract method 0x97973043.
//
// Solidity: function getSubmissionInfo(address _submissionID) view returns(uint8 status, uint64 submissionTime, uint64 index, bool registered, bool hasVouched, uint256 numberOfRequests)
func (_Poh *PohCallerSession) GetSubmissionInfo(_submissionID common.Address) (struct {
	Status           uint8
	SubmissionTime   uint64
	Index            uint64
	Registered       bool
	HasVouched       bool
	NumberOfRequests *big.Int
}, error) {
	return _Poh.Contract.GetSubmissionInfo(&_Poh.CallOpts, _submissionID)
}

// Governor is a free data retrieval call binding the contract method 0x0c340a24.
//
// Solidity: function governor() view returns(address)
func (_Poh *PohCaller) Governor(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Poh.contract.Call(opts, &out, "governor")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Governor is a free data retrieval call binding the contract method 0x0c340a24.
//
// Solidity: function governor() view returns(address)
func (_Poh *PohSession) Governor() (common.Address, error) {
	return _Poh.Contract.Governor(&_Poh.CallOpts)
}

// Governor is a free data retrieval call binding the contract method 0x0c340a24.
//
// Solidity: function governor() view returns(address)
func (_Poh *PohCallerSession) Governor() (common.Address, error) {
	return _Poh.Contract.Governor(&_Poh.CallOpts)
}

// IsRegistered is a free data retrieval call binding the contract method 0xc3c5a547.
//
// Solidity: function isRegistered(address _submissionID) view returns(bool)
func (_Poh *PohCaller) IsRegistered(opts *bind.CallOpts, _submissionID common.Address) (bool, error) {
	var out []interface{}
	err := _Poh.contract.Call(opts, &out, "isRegistered", _submissionID)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsRegistered is a free data retrieval call binding the contract method 0xc3c5a547.
//
// Solidity: function isRegistered(address _submissionID) view returns(bool)
func (_Poh *PohSession) IsRegistered(_submissionID common.Address) (bool, error) {
	return _Poh.Contract.IsRegistered(&_Poh.CallOpts, _submissionID)
}

// IsRegistered is a free data retrieval call binding the contract method 0xc3c5a547.
//
// Solidity: function isRegistered(address _submissionID) view returns(bool)
func (_Poh *PohCallerSession) IsRegistered(_submissionID common.Address) (bool, error) {
	return _Poh.Contract.IsRegistered(&_Poh.CallOpts, _submissionID)
}

// LoserStakeMultiplier is a free data retrieval call binding the contract method 0x1d512085.
//
// Solidity: function loserStakeMultiplier() view returns(uint256)
func (_Poh *PohCaller) LoserStakeMultiplier(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Poh.contract.Call(opts, &out, "loserStakeMultiplier")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LoserStakeMultiplier is a free data retrieval call binding the contract method 0x1d512085.
//
// Solidity: function loserStakeMultiplier() view returns(uint256)
func (_Poh *PohSession) LoserStakeMultiplier() (*big.Int, error) {
	return _Poh.Contract.LoserStakeMultiplier(&_Poh.CallOpts)
}

// LoserStakeMultiplier is a free data retrieval call binding the contract method 0x1d512085.
//
// Solidity: function loserStakeMultiplier() view returns(uint256)
func (_Poh *PohCallerSession) LoserStakeMultiplier() (*big.Int, error) {
	return _Poh.Contract.LoserStakeMultiplier(&_Poh.CallOpts)
}

// RenewalPeriodDuration is a free data retrieval call binding the contract method 0x876c63d4.
//
// Solidity: function renewalPeriodDuration() view returns(uint64)
func (_Poh *PohCaller) RenewalPeriodDuration(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Poh.contract.Call(opts, &out, "renewalPeriodDuration")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// RenewalPeriodDuration is a free data retrieval call binding the contract method 0x876c63d4.
//
// Solidity: function renewalPeriodDuration() view returns(uint64)
func (_Poh *PohSession) RenewalPeriodDuration() (uint64, error) {
	return _Poh.Contract.RenewalPeriodDuration(&_Poh.CallOpts)
}

// RenewalPeriodDuration is a free data retrieval call binding the contract method 0x876c63d4.
//
// Solidity: function renewalPeriodDuration() view returns(uint64)
func (_Poh *PohCallerSession) RenewalPeriodDuration() (uint64, error) {
	return _Poh.Contract.RenewalPeriodDuration(&_Poh.CallOpts)
}

// RequiredNumberOfVouches is a free data retrieval call binding the contract method 0x2d9489c6.
//
// Solidity: function requiredNumberOfVouches() view returns(uint64)
func (_Poh *PohCaller) RequiredNumberOfVouches(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Poh.contract.Call(opts, &out, "requiredNumberOfVouches")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// RequiredNumberOfVouches is a free data retrieval call binding the contract method 0x2d9489c6.
//
// Solidity: function requiredNumberOfVouches() view returns(uint64)
func (_Poh *PohSession) RequiredNumberOfVouches() (uint64, error) {
	return _Poh.Contract.RequiredNumberOfVouches(&_Poh.CallOpts)
}

// RequiredNumberOfVouches is a free data retrieval call binding the contract method 0x2d9489c6.
//
// Solidity: function requiredNumberOfVouches() view returns(uint64)
func (_Poh *PohCallerSession) RequiredNumberOfVouches() (uint64, error) {
	return _Poh.Contract.RequiredNumberOfVouches(&_Poh.CallOpts)
}

// SharedStakeMultiplier is a free data retrieval call binding the contract method 0x41658341.
//
// Solidity: function sharedStakeMultiplier() view returns(uint256)
func (_Poh *PohCaller) SharedStakeMultiplier(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Poh.contract.Call(opts, &out, "sharedStakeMultiplier")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SharedStakeMultiplier is a free data retrieval call binding the contract method 0x41658341.
//
// Solidity: function sharedStakeMultiplier() view returns(uint256)
func (_Poh *PohSession) SharedStakeMultiplier() (*big.Int, error) {
	return _Poh.Contract.SharedStakeMultiplier(&_Poh.CallOpts)
}

// SharedStakeMultiplier is a free data retrieval call binding the contract method 0x41658341.
//
// Solidity: function sharedStakeMultiplier() view returns(uint256)
func (_Poh *PohCallerSession) SharedStakeMultiplier() (*big.Int, error) {
	return _Poh.Contract.SharedStakeMultiplier(&_Poh.CallOpts)
}

// SubmissionBaseDeposit is a free data retrieval call binding the contract method 0xbb0b86ff.
//
// Solidity: function submissionBaseDeposit() view returns(uint256)
func (_Poh *PohCaller) SubmissionBaseDeposit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Poh.contract.Call(opts, &out, "submissionBaseDeposit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SubmissionBaseDeposit is a free data retrieval call binding the contract method 0xbb0b86ff.
//
// Solidity: function submissionBaseDeposit() view returns(uint256)
func (_Poh *PohSession) SubmissionBaseDeposit() (*big.Int, error) {
	return _Poh.Contract.SubmissionBaseDeposit(&_Poh.CallOpts)
}

// SubmissionBaseDeposit is a free data retrieval call binding the contract method 0xbb0b86ff.
//
// Solidity: function submissionBaseDeposit() view returns(uint256)
func (_Poh *PohCallerSession) SubmissionBaseDeposit() (*big.Int, error) {
	return _Poh.Contract.SubmissionBaseDeposit(&_Poh.CallOpts)
}

// SubmissionCounter is a free data retrieval call binding the contract method 0x76c23ff1.
//
// Solidity: function submissionCounter() view returns(uint256)
func (_Poh *PohCaller) SubmissionCounter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Poh.contract.Call(opts, &out, "submissionCounter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SubmissionCounter is a free data retrieval call binding the contract method 0x76c23ff1.
//
// Solidity: function submissionCounter() view returns(uint256)
func (_Poh *PohSession) SubmissionCounter() (*big.Int, error) {
	return _Poh.Contract.SubmissionCounter(&_Poh.CallOpts)
}

// SubmissionCounter is a free data retrieval call binding the contract method 0x76c23ff1.
//
// Solidity: function submissionCounter() view returns(uint256)
func (_Poh *PohCallerSession) SubmissionCounter() (*big.Int, error) {
	return _Poh.Contract.SubmissionCounter(&_Poh.CallOpts)
}

// SubmissionDuration is a free data retrieval call binding the contract method 0xf633c293.
//
// Solidity: function submissionDuration() view returns(uint64)
func (_Poh *PohCaller) SubmissionDuration(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Poh.contract.Call(opts, &out, "submissionDuration")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// SubmissionDuration is a free data retrieval call binding the contract method 0xf633c293.
//
// Solidity: function submissionDuration() view returns(uint64)
func (_Poh *PohSession) SubmissionDuration() (uint64, error) {
	return _Poh.Contract.SubmissionDuration(&_Poh.CallOpts)
}

// SubmissionDuration is a free data retrieval call binding the contract method 0xf633c293.
//
// Solidity: function submissionDuration() view returns(uint64)
func (_Poh *PohCallerSession) SubmissionDuration() (uint64, error) {
	return _Poh.Contract.SubmissionDuration(&_Poh.CallOpts)
}

// Vouches is a free data retrieval call binding the contract method 0x0b337be6.
//
// Solidity: function vouches(address , address ) view returns(bool)
func (_Poh *PohCaller) Vouches(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (bool, error) {
	var out []interface{}
	err := _Poh.contract.Call(opts, &out, "vouches", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Vouches is a free data retrieval call binding the contract method 0x0b337be6.
//
// Solidity: function vouches(address , address ) view returns(bool)
func (_Poh *PohSession) Vouches(arg0 common.Address, arg1 common.Address) (bool, error) {
	return _Poh.Contract.Vouches(&_Poh.CallOpts, arg0, arg1)
}

// Vouches is a free data retrieval call binding the contract method 0x0b337be6.
//
// Solidity: function vouches(address , address ) view returns(bool)
func (_Poh *PohCallerSession) Vouches(arg0 common.Address, arg1 common.Address) (bool, error) {
	return _Poh.Contract.Vouches(&_Poh.CallOpts, arg0, arg1)
}

// WinnerStakeMultiplier is a free data retrieval call binding the contract method 0x7b943383.
//
// Solidity: function winnerStakeMultiplier() view returns(uint256)
func (_Poh *PohCaller) WinnerStakeMultiplier(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Poh.contract.Call(opts, &out, "winnerStakeMultiplier")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WinnerStakeMultiplier is a free data retrieval call binding the contract method 0x7b943383.
//
// Solidity: function winnerStakeMultiplier() view returns(uint256)
func (_Poh *PohSession) WinnerStakeMultiplier() (*big.Int, error) {
	return _Poh.Contract.WinnerStakeMultiplier(&_Poh.CallOpts)
}

// WinnerStakeMultiplier is a free data retrieval call binding the contract method 0x7b943383.
//
// Solidity: function winnerStakeMultiplier() view returns(uint256)
func (_Poh *PohCallerSession) WinnerStakeMultiplier() (*big.Int, error) {
	return _Poh.Contract.WinnerStakeMultiplier(&_Poh.CallOpts)
}

// AddSubmission is a paid mutator transaction binding the contract method 0x4690d55a.
//
// Solidity: function addSubmission(string _evidence, string _name) payable returns()
func (_Poh *PohTransactor) AddSubmission(opts *bind.TransactOpts, _evidence string, _name string) (*types.Transaction, error) {
	return _Poh.contract.Transact(opts, "addSubmission", _evidence, _name)
}

// AddSubmission is a paid mutator transaction binding the contract method 0x4690d55a.
//
// Solidity: function addSubmission(string _evidence, string _name) payable returns()
func (_Poh *PohSession) AddSubmission(_evidence string, _name string) (*types.Transaction, error) {
	return _Poh.Contract.AddSubmission(&_Poh.TransactOpts, _evidence, _name)
}

// AddSubmission is a paid mutator transaction binding the contract method 0x4690d55a.
//
// Solidity: function addSubmission(string _evidence, string _name) payable returns()
func (_Poh *PohTransactorSession) AddSubmission(_evidence string, _name string) (*types.Transaction, error) {
	return _Poh.Contract.AddSubmission(&_Poh.TransactOpts, _evidence, _name)
}

// AddSubmissionManually is a paid mutator transaction binding the contract method 0x61b90541.
//
// Solidity: function addSubmissionManually(address[] _submissionIDs, string[] _evidence, string[] _names) returns()
func (_Poh *PohTransactor) AddSubmissionManually(opts *bind.TransactOpts, _submissionIDs []common.Address, _evidence []string, _names []string) (*types.Transaction, error) {
	return _Poh.contract.Transact(opts, "addSubmissionManually", _submissionIDs, _evidence, _names)
}

// AddSubmissionManually is a paid mutator transaction binding the contract method 0x61b90541.
//
// Solidity: function addSubmissionManually(address[] _submissionIDs, string[] _evidence, string[] _names) returns()
func (_Poh *PohSession) AddSubmissionManually(_submissionIDs []common.Address, _evidence []string, _names []string) (*types.Transaction, error) {
	return _Poh.Contract.AddSubmissionManually(&_Poh.TransactOpts, _submissionIDs, _evidence, _names)
}

// AddSubmissionManually is a paid mutator transaction binding the contract method 0x61b90541.
//
// Solidity: function addSubmissionManually(address[] _submissionIDs, string[] _evidence, string[] _names) returns()
func (_Poh *PohTransactorSession) AddSubmissionManually(_submissionIDs []common.Address, _evidence []string, _names []string) (*types.Transaction, error) {
	return _Poh.Contract.AddSubmissionManually(&_Poh.TransactOpts, _submissionIDs, _evidence, _names)
}

// AddVouch is a paid mutator transaction binding the contract method 0x32fe596f.
//
// Solidity: function addVouch(address _submissionID) returns()
func (_Poh *PohTransactor) AddVouch(opts *bind.TransactOpts, _submissionID common.Address) (*types.Transaction, error) {
	return _Poh.contract.Transact(opts, "addVouch", _submissionID)
}

// AddVouch is a paid mutator transaction binding the contract method 0x32fe596f.
//
// Solidity: function addVouch(address _submissionID) returns()
func (_Poh *PohSession) AddVouch(_submissionID common.Address) (*types.Transaction, error) {
	return _Poh.Contract.AddVouch(&_Poh.TransactOpts, _submissionID)
}

// AddVouch is a paid mutator transaction binding the contract method 0x32fe596f.
//
// Solidity: function addVouch(address _submissionID) returns()
func (_Poh *PohTransactorSession) AddVouch(_submissionID common.Address) (*types.Transaction, error) {
	return _Poh.Contract.AddVouch(&_Poh.TransactOpts, _submissionID)
}

// ChallengeRequest is a paid mutator transaction binding the contract method 0xf40e0aed.
//
// Solidity: function challengeRequest(address _submissionID, uint8 _reason, address _duplicateID, string _evidence) payable returns()
func (_Poh *PohTransactor) ChallengeRequest(opts *bind.TransactOpts, _submissionID common.Address, _reason uint8, _duplicateID common.Address, _evidence string) (*types.Transaction, error) {
	return _Poh.contract.Transact(opts, "challengeRequest", _submissionID, _reason, _duplicateID, _evidence)
}

// ChallengeRequest is a paid mutator transaction binding the contract method 0xf40e0aed.
//
// Solidity: function challengeRequest(address _submissionID, uint8 _reason, address _duplicateID, string _evidence) payable returns()
func (_Poh *PohSession) ChallengeRequest(_submissionID common.Address, _reason uint8, _duplicateID common.Address, _evidence string) (*types.Transaction, error) {
	return _Poh.Contract.ChallengeRequest(&_Poh.TransactOpts, _submissionID, _reason, _duplicateID, _evidence)
}

// ChallengeRequest is a paid mutator transaction binding the contract method 0xf40e0aed.
//
// Solidity: function challengeRequest(address _submissionID, uint8 _reason, address _duplicateID, string _evidence) payable returns()
func (_Poh *PohTransactorSession) ChallengeRequest(_submissionID common.Address, _reason uint8, _duplicateID common.Address, _evidence string) (*types.Transaction, error) {
	return _Poh.Contract.ChallengeRequest(&_Poh.TransactOpts, _submissionID, _reason, _duplicateID, _evidence)
}

// ChangeArbitrator is a paid mutator transaction binding the contract method 0xba7079ca.
//
// Solidity: function changeArbitrator(address _arbitrator, bytes _arbitratorExtraData) returns()
func (_Poh *PohTransactor) ChangeArbitrator(opts *bind.TransactOpts, _arbitrator common.Address, _arbitratorExtraData []byte) (*types.Transaction, error) {
	return _Poh.contract.Transact(opts, "changeArbitrator", _arbitrator, _arbitratorExtraData)
}

// ChangeArbitrator is a paid mutator transaction binding the contract method 0xba7079ca.
//
// Solidity: function changeArbitrator(address _arbitrator, bytes _arbitratorExtraData) returns()
func (_Poh *PohSession) ChangeArbitrator(_arbitrator common.Address, _arbitratorExtraData []byte) (*types.Transaction, error) {
	return _Poh.Contract.ChangeArbitrator(&_Poh.TransactOpts, _arbitrator, _arbitratorExtraData)
}

// ChangeArbitrator is a paid mutator transaction binding the contract method 0xba7079ca.
//
// Solidity: function changeArbitrator(address _arbitrator, bytes _arbitratorExtraData) returns()
func (_Poh *PohTransactorSession) ChangeArbitrator(_arbitrator common.Address, _arbitratorExtraData []byte) (*types.Transaction, error) {
	return _Poh.Contract.ChangeArbitrator(&_Poh.TransactOpts, _arbitrator, _arbitratorExtraData)
}

// ChangeDurations is a paid mutator transaction binding the contract method 0x26bafe5f.
//
// Solidity: function changeDurations(uint64 _submissionDuration, uint64 _renewalPeriodDuration, uint64 _challengePeriodDuration) returns()
func (_Poh *PohTransactor) ChangeDurations(opts *bind.TransactOpts, _submissionDuration uint64, _renewalPeriodDuration uint64, _challengePeriodDuration uint64) (*types.Transaction, error) {
	return _Poh.contract.Transact(opts, "changeDurations", _submissionDuration, _renewalPeriodDuration, _challengePeriodDuration)
}

// ChangeDurations is a paid mutator transaction binding the contract method 0x26bafe5f.
//
// Solidity: function changeDurations(uint64 _submissionDuration, uint64 _renewalPeriodDuration, uint64 _challengePeriodDuration) returns()
func (_Poh *PohSession) ChangeDurations(_submissionDuration uint64, _renewalPeriodDuration uint64, _challengePeriodDuration uint64) (*types.Transaction, error) {
	return _Poh.Contract.ChangeDurations(&_Poh.TransactOpts, _submissionDuration, _renewalPeriodDuration, _challengePeriodDuration)
}

// ChangeDurations is a paid mutator transaction binding the contract method 0x26bafe5f.
//
// Solidity: function changeDurations(uint64 _submissionDuration, uint64 _renewalPeriodDuration, uint64 _challengePeriodDuration) returns()
func (_Poh *PohTransactorSession) ChangeDurations(_submissionDuration uint64, _renewalPeriodDuration uint64, _challengePeriodDuration uint64) (*types.Transaction, error) {
	return _Poh.Contract.ChangeDurations(&_Poh.TransactOpts, _submissionDuration, _renewalPeriodDuration, _challengePeriodDuration)
}

// ChangeGovernor is a paid mutator transaction binding the contract method 0xe4c0aaf4.
//
// Solidity: function changeGovernor(address _governor) returns()
func (_Poh *PohTransactor) ChangeGovernor(opts *bind.TransactOpts, _governor common.Address) (*types.Transaction, error) {
	return _Poh.contract.Transact(opts, "changeGovernor", _governor)
}

// ChangeGovernor is a paid mutator transaction binding the contract method 0xe4c0aaf4.
//
// Solidity: function changeGovernor(address _governor) returns()
func (_Poh *PohSession) ChangeGovernor(_governor common.Address) (*types.Transaction, error) {
	return _Poh.Contract.ChangeGovernor(&_Poh.TransactOpts, _governor)
}

// ChangeGovernor is a paid mutator transaction binding the contract method 0xe4c0aaf4.
//
// Solidity: function changeGovernor(address _governor) returns()
func (_Poh *PohTransactorSession) ChangeGovernor(_governor common.Address) (*types.Transaction, error) {
	return _Poh.Contract.ChangeGovernor(&_Poh.TransactOpts, _governor)
}

// ChangeLoserStakeMultiplier is a paid mutator transaction binding the contract method 0x92239dff.
//
// Solidity: function changeLoserStakeMultiplier(uint256 _loserStakeMultiplier) returns()
func (_Poh *PohTransactor) ChangeLoserStakeMultiplier(opts *bind.TransactOpts, _loserStakeMultiplier *big.Int) (*types.Transaction, error) {
	return _Poh.contract.Transact(opts, "changeLoserStakeMultiplier", _loserStakeMultiplier)
}

// ChangeLoserStakeMultiplier is a paid mutator transaction binding the contract method 0x92239dff.
//
// Solidity: function changeLoserStakeMultiplier(uint256 _loserStakeMultiplier) returns()
func (_Poh *PohSession) ChangeLoserStakeMultiplier(_loserStakeMultiplier *big.Int) (*types.Transaction, error) {
	return _Poh.Contract.ChangeLoserStakeMultiplier(&_Poh.TransactOpts, _loserStakeMultiplier)
}

// ChangeLoserStakeMultiplier is a paid mutator transaction binding the contract method 0x92239dff.
//
// Solidity: function changeLoserStakeMultiplier(uint256 _loserStakeMultiplier) returns()
func (_Poh *PohTransactorSession) ChangeLoserStakeMultiplier(_loserStakeMultiplier *big.Int) (*types.Transaction, error) {
	return _Poh.Contract.ChangeLoserStakeMultiplier(&_Poh.TransactOpts, _loserStakeMultiplier)
}

// ChangeMetaEvidence is a paid mutator transaction binding the contract method 0xd706be31.
//
// Solidity: function changeMetaEvidence(string _registrationMetaEvidence, string _clearingMetaEvidence) returns()
func (_Poh *PohTransactor) ChangeMetaEvidence(opts *bind.TransactOpts, _registrationMetaEvidence string, _clearingMetaEvidence string) (*types.Transaction, error) {
	return _Poh.contract.Transact(opts, "changeMetaEvidence", _registrationMetaEvidence, _clearingMetaEvidence)
}

// ChangeMetaEvidence is a paid mutator transaction binding the contract method 0xd706be31.
//
// Solidity: function changeMetaEvidence(string _registrationMetaEvidence, string _clearingMetaEvidence) returns()
func (_Poh *PohSession) ChangeMetaEvidence(_registrationMetaEvidence string, _clearingMetaEvidence string) (*types.Transaction, error) {
	return _Poh.Contract.ChangeMetaEvidence(&_Poh.TransactOpts, _registrationMetaEvidence, _clearingMetaEvidence)
}

// ChangeMetaEvidence is a paid mutator transaction binding the contract method 0xd706be31.
//
// Solidity: function changeMetaEvidence(string _registrationMetaEvidence, string _clearingMetaEvidence) returns()
func (_Poh *PohTransactorSession) ChangeMetaEvidence(_registrationMetaEvidence string, _clearingMetaEvidence string) (*types.Transaction, error) {
	return _Poh.Contract.ChangeMetaEvidence(&_Poh.TransactOpts, _registrationMetaEvidence, _clearingMetaEvidence)
}

// ChangeRequiredNumberOfVouches is a paid mutator transaction binding the contract method 0xf65ab1be.
//
// Solidity: function changeRequiredNumberOfVouches(uint64 _requiredNumberOfVouches) returns()
func (_Poh *PohTransactor) ChangeRequiredNumberOfVouches(opts *bind.TransactOpts, _requiredNumberOfVouches uint64) (*types.Transaction, error) {
	return _Poh.contract.Transact(opts, "changeRequiredNumberOfVouches", _requiredNumberOfVouches)
}

// ChangeRequiredNumberOfVouches is a paid mutator transaction binding the contract method 0xf65ab1be.
//
// Solidity: function changeRequiredNumberOfVouches(uint64 _requiredNumberOfVouches) returns()
func (_Poh *PohSession) ChangeRequiredNumberOfVouches(_requiredNumberOfVouches uint64) (*types.Transaction, error) {
	return _Poh.Contract.ChangeRequiredNumberOfVouches(&_Poh.TransactOpts, _requiredNumberOfVouches)
}

// ChangeRequiredNumberOfVouches is a paid mutator transaction binding the contract method 0xf65ab1be.
//
// Solidity: function changeRequiredNumberOfVouches(uint64 _requiredNumberOfVouches) returns()
func (_Poh *PohTransactorSession) ChangeRequiredNumberOfVouches(_requiredNumberOfVouches uint64) (*types.Transaction, error) {
	return _Poh.Contract.ChangeRequiredNumberOfVouches(&_Poh.TransactOpts, _requiredNumberOfVouches)
}

// ChangeSharedStakeMultiplier is a paid mutator transaction binding the contract method 0x12ce3525.
//
// Solidity: function changeSharedStakeMultiplier(uint256 _sharedStakeMultiplier) returns()
func (_Poh *PohTransactor) ChangeSharedStakeMultiplier(opts *bind.TransactOpts, _sharedStakeMultiplier *big.Int) (*types.Transaction, error) {
	return _Poh.contract.Transact(opts, "changeSharedStakeMultiplier", _sharedStakeMultiplier)
}

// ChangeSharedStakeMultiplier is a paid mutator transaction binding the contract method 0x12ce3525.
//
// Solidity: function changeSharedStakeMultiplier(uint256 _sharedStakeMultiplier) returns()
func (_Poh *PohSession) ChangeSharedStakeMultiplier(_sharedStakeMultiplier *big.Int) (*types.Transaction, error) {
	return _Poh.Contract.ChangeSharedStakeMultiplier(&_Poh.TransactOpts, _sharedStakeMultiplier)
}

// ChangeSharedStakeMultiplier is a paid mutator transaction binding the contract method 0x12ce3525.
//
// Solidity: function changeSharedStakeMultiplier(uint256 _sharedStakeMultiplier) returns()
func (_Poh *PohTransactorSession) ChangeSharedStakeMultiplier(_sharedStakeMultiplier *big.Int) (*types.Transaction, error) {
	return _Poh.Contract.ChangeSharedStakeMultiplier(&_Poh.TransactOpts, _sharedStakeMultiplier)
}

// ChangeStateToPending is a paid mutator transaction binding the contract method 0xb4dfe93d.
//
// Solidity: function changeStateToPending(address _submissionID, address[] _vouches, bytes[] _signatures, uint256[] _expirationTimestamps) returns()
func (_Poh *PohTransactor) ChangeStateToPending(opts *bind.TransactOpts, _submissionID common.Address, _vouches []common.Address, _signatures [][]byte, _expirationTimestamps []*big.Int) (*types.Transaction, error) {
	return _Poh.contract.Transact(opts, "changeStateToPending", _submissionID, _vouches, _signatures, _expirationTimestamps)
}

// ChangeStateToPending is a paid mutator transaction binding the contract method 0xb4dfe93d.
//
// Solidity: function changeStateToPending(address _submissionID, address[] _vouches, bytes[] _signatures, uint256[] _expirationTimestamps) returns()
func (_Poh *PohSession) ChangeStateToPending(_submissionID common.Address, _vouches []common.Address, _signatures [][]byte, _expirationTimestamps []*big.Int) (*types.Transaction, error) {
	return _Poh.Contract.ChangeStateToPending(&_Poh.TransactOpts, _submissionID, _vouches, _signatures, _expirationTimestamps)
}

// ChangeStateToPending is a paid mutator transaction binding the contract method 0xb4dfe93d.
//
// Solidity: function changeStateToPending(address _submissionID, address[] _vouches, bytes[] _signatures, uint256[] _expirationTimestamps) returns()
func (_Poh *PohTransactorSession) ChangeStateToPending(_submissionID common.Address, _vouches []common.Address, _signatures [][]byte, _expirationTimestamps []*big.Int) (*types.Transaction, error) {
	return _Poh.Contract.ChangeStateToPending(&_Poh.TransactOpts, _submissionID, _vouches, _signatures, _expirationTimestamps)
}

// ChangeSubmissionBaseDeposit is a paid mutator transaction binding the contract method 0x33e5d047.
//
// Solidity: function changeSubmissionBaseDeposit(uint256 _submissionBaseDeposit) returns()
func (_Poh *PohTransactor) ChangeSubmissionBaseDeposit(opts *bind.TransactOpts, _submissionBaseDeposit *big.Int) (*types.Transaction, error) {
	return _Poh.contract.Transact(opts, "changeSubmissionBaseDeposit", _submissionBaseDeposit)
}

// ChangeSubmissionBaseDeposit is a paid mutator transaction binding the contract method 0x33e5d047.
//
// Solidity: function changeSubmissionBaseDeposit(uint256 _submissionBaseDeposit) returns()
func (_Poh *PohSession) ChangeSubmissionBaseDeposit(_submissionBaseDeposit *big.Int) (*types.Transaction, error) {
	return _Poh.Contract.ChangeSubmissionBaseDeposit(&_Poh.TransactOpts, _submissionBaseDeposit)
}

// ChangeSubmissionBaseDeposit is a paid mutator transaction binding the contract method 0x33e5d047.
//
// Solidity: function changeSubmissionBaseDeposit(uint256 _submissionBaseDeposit) returns()
func (_Poh *PohTransactorSession) ChangeSubmissionBaseDeposit(_submissionBaseDeposit *big.Int) (*types.Transaction, error) {
	return _Poh.Contract.ChangeSubmissionBaseDeposit(&_Poh.TransactOpts, _submissionBaseDeposit)
}

// ChangeWinnerStakeMultiplier is a paid mutator transaction binding the contract method 0xadc7faba.
//
// Solidity: function changeWinnerStakeMultiplier(uint256 _winnerStakeMultiplier) returns()
func (_Poh *PohTransactor) ChangeWinnerStakeMultiplier(opts *bind.TransactOpts, _winnerStakeMultiplier *big.Int) (*types.Transaction, error) {
	return _Poh.contract.Transact(opts, "changeWinnerStakeMultiplier", _winnerStakeMultiplier)
}

// ChangeWinnerStakeMultiplier is a paid mutator transaction binding the contract method 0xadc7faba.
//
// Solidity: function changeWinnerStakeMultiplier(uint256 _winnerStakeMultiplier) returns()
func (_Poh *PohSession) ChangeWinnerStakeMultiplier(_winnerStakeMultiplier *big.Int) (*types.Transaction, error) {
	return _Poh.Contract.ChangeWinnerStakeMultiplier(&_Poh.TransactOpts, _winnerStakeMultiplier)
}

// ChangeWinnerStakeMultiplier is a paid mutator transaction binding the contract method 0xadc7faba.
//
// Solidity: function changeWinnerStakeMultiplier(uint256 _winnerStakeMultiplier) returns()
func (_Poh *PohTransactorSession) ChangeWinnerStakeMultiplier(_winnerStakeMultiplier *big.Int) (*types.Transaction, error) {
	return _Poh.Contract.ChangeWinnerStakeMultiplier(&_Poh.TransactOpts, _winnerStakeMultiplier)
}

// ExecuteRequest is a paid mutator transaction binding the contract method 0x6e98762d.
//
// Solidity: function executeRequest(address _submissionID) returns()
func (_Poh *PohTransactor) ExecuteRequest(opts *bind.TransactOpts, _submissionID common.Address) (*types.Transaction, error) {
	return _Poh.contract.Transact(opts, "executeRequest", _submissionID)
}

// ExecuteRequest is a paid mutator transaction binding the contract method 0x6e98762d.
//
// Solidity: function executeRequest(address _submissionID) returns()
func (_Poh *PohSession) ExecuteRequest(_submissionID common.Address) (*types.Transaction, error) {
	return _Poh.Contract.ExecuteRequest(&_Poh.TransactOpts, _submissionID)
}

// ExecuteRequest is a paid mutator transaction binding the contract method 0x6e98762d.
//
// Solidity: function executeRequest(address _submissionID) returns()
func (_Poh *PohTransactorSession) ExecuteRequest(_submissionID common.Address) (*types.Transaction, error) {
	return _Poh.Contract.ExecuteRequest(&_Poh.TransactOpts, _submissionID)
}

// FundAppeal is a paid mutator transaction binding the contract method 0xd7e9f178.
//
// Solidity: function fundAppeal(address _submissionID, uint256 _challengeID, uint8 _side) payable returns()
func (_Poh *PohTransactor) FundAppeal(opts *bind.TransactOpts, _submissionID common.Address, _challengeID *big.Int, _side uint8) (*types.Transaction, error) {
	return _Poh.contract.Transact(opts, "fundAppeal", _submissionID, _challengeID, _side)
}

// FundAppeal is a paid mutator transaction binding the contract method 0xd7e9f178.
//
// Solidity: function fundAppeal(address _submissionID, uint256 _challengeID, uint8 _side) payable returns()
func (_Poh *PohSession) FundAppeal(_submissionID common.Address, _challengeID *big.Int, _side uint8) (*types.Transaction, error) {
	return _Poh.Contract.FundAppeal(&_Poh.TransactOpts, _submissionID, _challengeID, _side)
}

// FundAppeal is a paid mutator transaction binding the contract method 0xd7e9f178.
//
// Solidity: function fundAppeal(address _submissionID, uint256 _challengeID, uint8 _side) payable returns()
func (_Poh *PohTransactorSession) FundAppeal(_submissionID common.Address, _challengeID *big.Int, _side uint8) (*types.Transaction, error) {
	return _Poh.Contract.FundAppeal(&_Poh.TransactOpts, _submissionID, _challengeID, _side)
}

// FundSubmission is a paid mutator transaction binding the contract method 0xa27456bb.
//
// Solidity: function fundSubmission(address _submissionID) payable returns()
func (_Poh *PohTransactor) FundSubmission(opts *bind.TransactOpts, _submissionID common.Address) (*types.Transaction, error) {
	return _Poh.contract.Transact(opts, "fundSubmission", _submissionID)
}

// FundSubmission is a paid mutator transaction binding the contract method 0xa27456bb.
//
// Solidity: function fundSubmission(address _submissionID) payable returns()
func (_Poh *PohSession) FundSubmission(_submissionID common.Address) (*types.Transaction, error) {
	return _Poh.Contract.FundSubmission(&_Poh.TransactOpts, _submissionID)
}

// FundSubmission is a paid mutator transaction binding the contract method 0xa27456bb.
//
// Solidity: function fundSubmission(address _submissionID) payable returns()
func (_Poh *PohTransactorSession) FundSubmission(_submissionID common.Address) (*types.Transaction, error) {
	return _Poh.Contract.FundSubmission(&_Poh.TransactOpts, _submissionID)
}

// ProcessVouches is a paid mutator transaction binding the contract method 0x649a08bf.
//
// Solidity: function processVouches(address _submissionID, uint256 _requestID, uint256 _iterations) returns()
func (_Poh *PohTransactor) ProcessVouches(opts *bind.TransactOpts, _submissionID common.Address, _requestID *big.Int, _iterations *big.Int) (*types.Transaction, error) {
	return _Poh.contract.Transact(opts, "processVouches", _submissionID, _requestID, _iterations)
}

// ProcessVouches is a paid mutator transaction binding the contract method 0x649a08bf.
//
// Solidity: function processVouches(address _submissionID, uint256 _requestID, uint256 _iterations) returns()
func (_Poh *PohSession) ProcessVouches(_submissionID common.Address, _requestID *big.Int, _iterations *big.Int) (*types.Transaction, error) {
	return _Poh.Contract.ProcessVouches(&_Poh.TransactOpts, _submissionID, _requestID, _iterations)
}

// ProcessVouches is a paid mutator transaction binding the contract method 0x649a08bf.
//
// Solidity: function processVouches(address _submissionID, uint256 _requestID, uint256 _iterations) returns()
func (_Poh *PohTransactorSession) ProcessVouches(_submissionID common.Address, _requestID *big.Int, _iterations *big.Int) (*types.Transaction, error) {
	return _Poh.Contract.ProcessVouches(&_Poh.TransactOpts, _submissionID, _requestID, _iterations)
}

// ReapplySubmission is a paid mutator transaction binding the contract method 0x5a9ef341.
//
// Solidity: function reapplySubmission(string _evidence, string _name) payable returns()
func (_Poh *PohTransactor) ReapplySubmission(opts *bind.TransactOpts, _evidence string, _name string) (*types.Transaction, error) {
	return _Poh.contract.Transact(opts, "reapplySubmission", _evidence, _name)
}

// ReapplySubmission is a paid mutator transaction binding the contract method 0x5a9ef341.
//
// Solidity: function reapplySubmission(string _evidence, string _name) payable returns()
func (_Poh *PohSession) ReapplySubmission(_evidence string, _name string) (*types.Transaction, error) {
	return _Poh.Contract.ReapplySubmission(&_Poh.TransactOpts, _evidence, _name)
}

// ReapplySubmission is a paid mutator transaction binding the contract method 0x5a9ef341.
//
// Solidity: function reapplySubmission(string _evidence, string _name) payable returns()
func (_Poh *PohTransactorSession) ReapplySubmission(_evidence string, _name string) (*types.Transaction, error) {
	return _Poh.Contract.ReapplySubmission(&_Poh.TransactOpts, _evidence, _name)
}

// RemoveSubmission is a paid mutator transaction binding the contract method 0xf4934cdb.
//
// Solidity: function removeSubmission(address _submissionID, string _evidence) payable returns()
func (_Poh *PohTransactor) RemoveSubmission(opts *bind.TransactOpts, _submissionID common.Address, _evidence string) (*types.Transaction, error) {
	return _Poh.contract.Transact(opts, "removeSubmission", _submissionID, _evidence)
}

// RemoveSubmission is a paid mutator transaction binding the contract method 0xf4934cdb.
//
// Solidity: function removeSubmission(address _submissionID, string _evidence) payable returns()
func (_Poh *PohSession) RemoveSubmission(_submissionID common.Address, _evidence string) (*types.Transaction, error) {
	return _Poh.Contract.RemoveSubmission(&_Poh.TransactOpts, _submissionID, _evidence)
}

// RemoveSubmission is a paid mutator transaction binding the contract method 0xf4934cdb.
//
// Solidity: function removeSubmission(address _submissionID, string _evidence) payable returns()
func (_Poh *PohTransactorSession) RemoveSubmission(_submissionID common.Address, _evidence string) (*types.Transaction, error) {
	return _Poh.Contract.RemoveSubmission(&_Poh.TransactOpts, _submissionID, _evidence)
}

// RemoveSubmissionManually is a paid mutator transaction binding the contract method 0xa6c6ecc9.
//
// Solidity: function removeSubmissionManually(address _submissionID) returns()
func (_Poh *PohTransactor) RemoveSubmissionManually(opts *bind.TransactOpts, _submissionID common.Address) (*types.Transaction, error) {
	return _Poh.contract.Transact(opts, "removeSubmissionManually", _submissionID)
}

// RemoveSubmissionManually is a paid mutator transaction binding the contract method 0xa6c6ecc9.
//
// Solidity: function removeSubmissionManually(address _submissionID) returns()
func (_Poh *PohSession) RemoveSubmissionManually(_submissionID common.Address) (*types.Transaction, error) {
	return _Poh.Contract.RemoveSubmissionManually(&_Poh.TransactOpts, _submissionID)
}

// RemoveSubmissionManually is a paid mutator transaction binding the contract method 0xa6c6ecc9.
//
// Solidity: function removeSubmissionManually(address _submissionID) returns()
func (_Poh *PohTransactorSession) RemoveSubmissionManually(_submissionID common.Address) (*types.Transaction, error) {
	return _Poh.Contract.RemoveSubmissionManually(&_Poh.TransactOpts, _submissionID)
}

// RemoveVouch is a paid mutator transaction binding the contract method 0x84d1c91a.
//
// Solidity: function removeVouch(address _submissionID) returns()
func (_Poh *PohTransactor) RemoveVouch(opts *bind.TransactOpts, _submissionID common.Address) (*types.Transaction, error) {
	return _Poh.contract.Transact(opts, "removeVouch", _submissionID)
}

// RemoveVouch is a paid mutator transaction binding the contract method 0x84d1c91a.
//
// Solidity: function removeVouch(address _submissionID) returns()
func (_Poh *PohSession) RemoveVouch(_submissionID common.Address) (*types.Transaction, error) {
	return _Poh.Contract.RemoveVouch(&_Poh.TransactOpts, _submissionID)
}

// RemoveVouch is a paid mutator transaction binding the contract method 0x84d1c91a.
//
// Solidity: function removeVouch(address _submissionID) returns()
func (_Poh *PohTransactorSession) RemoveVouch(_submissionID common.Address) (*types.Transaction, error) {
	return _Poh.Contract.RemoveVouch(&_Poh.TransactOpts, _submissionID)
}

// Rule is a paid mutator transaction binding the contract method 0x311a6c56.
//
// Solidity: function rule(uint256 _disputeID, uint256 _ruling) returns()
func (_Poh *PohTransactor) Rule(opts *bind.TransactOpts, _disputeID *big.Int, _ruling *big.Int) (*types.Transaction, error) {
	return _Poh.contract.Transact(opts, "rule", _disputeID, _ruling)
}

// Rule is a paid mutator transaction binding the contract method 0x311a6c56.
//
// Solidity: function rule(uint256 _disputeID, uint256 _ruling) returns()
func (_Poh *PohSession) Rule(_disputeID *big.Int, _ruling *big.Int) (*types.Transaction, error) {
	return _Poh.Contract.Rule(&_Poh.TransactOpts, _disputeID, _ruling)
}

// Rule is a paid mutator transaction binding the contract method 0x311a6c56.
//
// Solidity: function rule(uint256 _disputeID, uint256 _ruling) returns()
func (_Poh *PohTransactorSession) Rule(_disputeID *big.Int, _ruling *big.Int) (*types.Transaction, error) {
	return _Poh.Contract.Rule(&_Poh.TransactOpts, _disputeID, _ruling)
}

// SubmitEvidence is a paid mutator transaction binding the contract method 0x5bb5e55b.
//
// Solidity: function submitEvidence(address _submissionID, string _evidence) returns()
func (_Poh *PohTransactor) SubmitEvidence(opts *bind.TransactOpts, _submissionID common.Address, _evidence string) (*types.Transaction, error) {
	return _Poh.contract.Transact(opts, "submitEvidence", _submissionID, _evidence)
}

// SubmitEvidence is a paid mutator transaction binding the contract method 0x5bb5e55b.
//
// Solidity: function submitEvidence(address _submissionID, string _evidence) returns()
func (_Poh *PohSession) SubmitEvidence(_submissionID common.Address, _evidence string) (*types.Transaction, error) {
	return _Poh.Contract.SubmitEvidence(&_Poh.TransactOpts, _submissionID, _evidence)
}

// SubmitEvidence is a paid mutator transaction binding the contract method 0x5bb5e55b.
//
// Solidity: function submitEvidence(address _submissionID, string _evidence) returns()
func (_Poh *PohTransactorSession) SubmitEvidence(_submissionID common.Address, _evidence string) (*types.Transaction, error) {
	return _Poh.Contract.SubmitEvidence(&_Poh.TransactOpts, _submissionID, _evidence)
}

// WithdrawFeesAndRewards is a paid mutator transaction binding the contract method 0x9a72e0b3.
//
// Solidity: function withdrawFeesAndRewards(address _beneficiary, address _submissionID, uint256 _requestID, uint256 _challengeID, uint256 _round) returns()
func (_Poh *PohTransactor) WithdrawFeesAndRewards(opts *bind.TransactOpts, _beneficiary common.Address, _submissionID common.Address, _requestID *big.Int, _challengeID *big.Int, _round *big.Int) (*types.Transaction, error) {
	return _Poh.contract.Transact(opts, "withdrawFeesAndRewards", _beneficiary, _submissionID, _requestID, _challengeID, _round)
}

// WithdrawFeesAndRewards is a paid mutator transaction binding the contract method 0x9a72e0b3.
//
// Solidity: function withdrawFeesAndRewards(address _beneficiary, address _submissionID, uint256 _requestID, uint256 _challengeID, uint256 _round) returns()
func (_Poh *PohSession) WithdrawFeesAndRewards(_beneficiary common.Address, _submissionID common.Address, _requestID *big.Int, _challengeID *big.Int, _round *big.Int) (*types.Transaction, error) {
	return _Poh.Contract.WithdrawFeesAndRewards(&_Poh.TransactOpts, _beneficiary, _submissionID, _requestID, _challengeID, _round)
}

// WithdrawFeesAndRewards is a paid mutator transaction binding the contract method 0x9a72e0b3.
//
// Solidity: function withdrawFeesAndRewards(address _beneficiary, address _submissionID, uint256 _requestID, uint256 _challengeID, uint256 _round) returns()
func (_Poh *PohTransactorSession) WithdrawFeesAndRewards(_beneficiary common.Address, _submissionID common.Address, _requestID *big.Int, _challengeID *big.Int, _round *big.Int) (*types.Transaction, error) {
	return _Poh.Contract.WithdrawFeesAndRewards(&_Poh.TransactOpts, _beneficiary, _submissionID, _requestID, _challengeID, _round)
}

// WithdrawSubmission is a paid mutator transaction binding the contract method 0xce52b9f4.
//
// Solidity: function withdrawSubmission() returns()
func (_Poh *PohTransactor) WithdrawSubmission(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Poh.contract.Transact(opts, "withdrawSubmission")
}

// WithdrawSubmission is a paid mutator transaction binding the contract method 0xce52b9f4.
//
// Solidity: function withdrawSubmission() returns()
func (_Poh *PohSession) WithdrawSubmission() (*types.Transaction, error) {
	return _Poh.Contract.WithdrawSubmission(&_Poh.TransactOpts)
}

// WithdrawSubmission is a paid mutator transaction binding the contract method 0xce52b9f4.
//
// Solidity: function withdrawSubmission() returns()
func (_Poh *PohTransactorSession) WithdrawSubmission() (*types.Transaction, error) {
	return _Poh.Contract.WithdrawSubmission(&_Poh.TransactOpts)
}

// PohAddSubmissionIterator is returned from FilterAddSubmission and is used to iterate over the raw logs and unpacked data for AddSubmission events raised by the Poh contract.
type PohAddSubmissionIterator struct {
	Event *PohAddSubmission // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PohAddSubmissionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PohAddSubmission)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PohAddSubmission)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PohAddSubmissionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PohAddSubmissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PohAddSubmission represents a AddSubmission event raised by the Poh contract.
type PohAddSubmission struct {
	SubmissionID common.Address
	RequestID    *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterAddSubmission is a free log retrieval operation binding the contract event 0x803727a67d35270dc2c090dc4f9cba1f9818a7200e65c2087eca187851fd6b19.
//
// Solidity: event AddSubmission(address indexed _submissionID, uint256 _requestID)
func (_Poh *PohFilterer) FilterAddSubmission(opts *bind.FilterOpts, _submissionID []common.Address) (*PohAddSubmissionIterator, error) {

	var _submissionIDRule []interface{}
	for _, _submissionIDItem := range _submissionID {
		_submissionIDRule = append(_submissionIDRule, _submissionIDItem)
	}

	logs, sub, err := _Poh.contract.FilterLogs(opts, "AddSubmission", _submissionIDRule)
	if err != nil {
		return nil, err
	}
	return &PohAddSubmissionIterator{contract: _Poh.contract, event: "AddSubmission", logs: logs, sub: sub}, nil
}

// WatchAddSubmission is a free log subscription operation binding the contract event 0x803727a67d35270dc2c090dc4f9cba1f9818a7200e65c2087eca187851fd6b19.
//
// Solidity: event AddSubmission(address indexed _submissionID, uint256 _requestID)
func (_Poh *PohFilterer) WatchAddSubmission(opts *bind.WatchOpts, sink chan<- *PohAddSubmission, _submissionID []common.Address) (event.Subscription, error) {

	var _submissionIDRule []interface{}
	for _, _submissionIDItem := range _submissionID {
		_submissionIDRule = append(_submissionIDRule, _submissionIDItem)
	}

	logs, sub, err := _Poh.contract.WatchLogs(opts, "AddSubmission", _submissionIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PohAddSubmission)
				if err := _Poh.contract.UnpackLog(event, "AddSubmission", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAddSubmission is a log parse operation binding the contract event 0x803727a67d35270dc2c090dc4f9cba1f9818a7200e65c2087eca187851fd6b19.
//
// Solidity: event AddSubmission(address indexed _submissionID, uint256 _requestID)
func (_Poh *PohFilterer) ParseAddSubmission(log types.Log) (*PohAddSubmission, error) {
	event := new(PohAddSubmission)
	if err := _Poh.contract.UnpackLog(event, "AddSubmission", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PohAppealContributionIterator is returned from FilterAppealContribution and is used to iterate over the raw logs and unpacked data for AppealContribution events raised by the Poh contract.
type PohAppealContributionIterator struct {
	Event *PohAppealContribution // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PohAppealContributionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PohAppealContribution)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PohAppealContribution)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PohAppealContributionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PohAppealContributionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PohAppealContribution represents a AppealContribution event raised by the Poh contract.
type PohAppealContribution struct {
	SubmissionID common.Address
	ChallengeID  *big.Int
	Party        uint8
	Contributor  common.Address
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterAppealContribution is a free log retrieval operation binding the contract event 0x9294febeba62e3f0e89187b59ba1235acc7fdbdebf15d6ede13c7b43adfa66dc.
//
// Solidity: event AppealContribution(address indexed _submissionID, uint256 indexed _challengeID, uint8 _party, address indexed _contributor, uint256 _amount)
func (_Poh *PohFilterer) FilterAppealContribution(opts *bind.FilterOpts, _submissionID []common.Address, _challengeID []*big.Int, _contributor []common.Address) (*PohAppealContributionIterator, error) {

	var _submissionIDRule []interface{}
	for _, _submissionIDItem := range _submissionID {
		_submissionIDRule = append(_submissionIDRule, _submissionIDItem)
	}
	var _challengeIDRule []interface{}
	for _, _challengeIDItem := range _challengeID {
		_challengeIDRule = append(_challengeIDRule, _challengeIDItem)
	}

	var _contributorRule []interface{}
	for _, _contributorItem := range _contributor {
		_contributorRule = append(_contributorRule, _contributorItem)
	}

	logs, sub, err := _Poh.contract.FilterLogs(opts, "AppealContribution", _submissionIDRule, _challengeIDRule, _contributorRule)
	if err != nil {
		return nil, err
	}
	return &PohAppealContributionIterator{contract: _Poh.contract, event: "AppealContribution", logs: logs, sub: sub}, nil
}

// WatchAppealContribution is a free log subscription operation binding the contract event 0x9294febeba62e3f0e89187b59ba1235acc7fdbdebf15d6ede13c7b43adfa66dc.
//
// Solidity: event AppealContribution(address indexed _submissionID, uint256 indexed _challengeID, uint8 _party, address indexed _contributor, uint256 _amount)
func (_Poh *PohFilterer) WatchAppealContribution(opts *bind.WatchOpts, sink chan<- *PohAppealContribution, _submissionID []common.Address, _challengeID []*big.Int, _contributor []common.Address) (event.Subscription, error) {

	var _submissionIDRule []interface{}
	for _, _submissionIDItem := range _submissionID {
		_submissionIDRule = append(_submissionIDRule, _submissionIDItem)
	}
	var _challengeIDRule []interface{}
	for _, _challengeIDItem := range _challengeID {
		_challengeIDRule = append(_challengeIDRule, _challengeIDItem)
	}

	var _contributorRule []interface{}
	for _, _contributorItem := range _contributor {
		_contributorRule = append(_contributorRule, _contributorItem)
	}

	logs, sub, err := _Poh.contract.WatchLogs(opts, "AppealContribution", _submissionIDRule, _challengeIDRule, _contributorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PohAppealContribution)
				if err := _Poh.contract.UnpackLog(event, "AppealContribution", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAppealContribution is a log parse operation binding the contract event 0x9294febeba62e3f0e89187b59ba1235acc7fdbdebf15d6ede13c7b43adfa66dc.
//
// Solidity: event AppealContribution(address indexed _submissionID, uint256 indexed _challengeID, uint8 _party, address indexed _contributor, uint256 _amount)
func (_Poh *PohFilterer) ParseAppealContribution(log types.Log) (*PohAppealContribution, error) {
	event := new(PohAppealContribution)
	if err := _Poh.contract.UnpackLog(event, "AppealContribution", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PohArbitratorCompleteIterator is returned from FilterArbitratorComplete and is used to iterate over the raw logs and unpacked data for ArbitratorComplete events raised by the Poh contract.
type PohArbitratorCompleteIterator struct {
	Event *PohArbitratorComplete // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PohArbitratorCompleteIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PohArbitratorComplete)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PohArbitratorComplete)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PohArbitratorCompleteIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PohArbitratorCompleteIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PohArbitratorComplete represents a ArbitratorComplete event raised by the Poh contract.
type PohArbitratorComplete struct {
	Arbitrator              common.Address
	Governor                common.Address
	SubmissionBaseDeposit   *big.Int
	SubmissionDuration      *big.Int
	ChallengePeriodDuration *big.Int
	RequiredNumberOfVouches *big.Int
	SharedStakeMultiplier   *big.Int
	WinnerStakeMultiplier   *big.Int
	LoserStakeMultiplier    *big.Int
	Raw                     types.Log // Blockchain specific contextual infos
}

// FilterArbitratorComplete is a free log retrieval operation binding the contract event 0xe313403816674ac13a5fd6524fafc793584c317c3c947825b71a753b0c155b20.
//
// Solidity: event ArbitratorComplete(address _arbitrator, address indexed _governor, uint256 _submissionBaseDeposit, uint256 _submissionDuration, uint256 _challengePeriodDuration, uint256 _requiredNumberOfVouches, uint256 _sharedStakeMultiplier, uint256 _winnerStakeMultiplier, uint256 _loserStakeMultiplier)
func (_Poh *PohFilterer) FilterArbitratorComplete(opts *bind.FilterOpts, _governor []common.Address) (*PohArbitratorCompleteIterator, error) {

	var _governorRule []interface{}
	for _, _governorItem := range _governor {
		_governorRule = append(_governorRule, _governorItem)
	}

	logs, sub, err := _Poh.contract.FilterLogs(opts, "ArbitratorComplete", _governorRule)
	if err != nil {
		return nil, err
	}
	return &PohArbitratorCompleteIterator{contract: _Poh.contract, event: "ArbitratorComplete", logs: logs, sub: sub}, nil
}

// WatchArbitratorComplete is a free log subscription operation binding the contract event 0xe313403816674ac13a5fd6524fafc793584c317c3c947825b71a753b0c155b20.
//
// Solidity: event ArbitratorComplete(address _arbitrator, address indexed _governor, uint256 _submissionBaseDeposit, uint256 _submissionDuration, uint256 _challengePeriodDuration, uint256 _requiredNumberOfVouches, uint256 _sharedStakeMultiplier, uint256 _winnerStakeMultiplier, uint256 _loserStakeMultiplier)
func (_Poh *PohFilterer) WatchArbitratorComplete(opts *bind.WatchOpts, sink chan<- *PohArbitratorComplete, _governor []common.Address) (event.Subscription, error) {

	var _governorRule []interface{}
	for _, _governorItem := range _governor {
		_governorRule = append(_governorRule, _governorItem)
	}

	logs, sub, err := _Poh.contract.WatchLogs(opts, "ArbitratorComplete", _governorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PohArbitratorComplete)
				if err := _Poh.contract.UnpackLog(event, "ArbitratorComplete", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseArbitratorComplete is a log parse operation binding the contract event 0xe313403816674ac13a5fd6524fafc793584c317c3c947825b71a753b0c155b20.
//
// Solidity: event ArbitratorComplete(address _arbitrator, address indexed _governor, uint256 _submissionBaseDeposit, uint256 _submissionDuration, uint256 _challengePeriodDuration, uint256 _requiredNumberOfVouches, uint256 _sharedStakeMultiplier, uint256 _winnerStakeMultiplier, uint256 _loserStakeMultiplier)
func (_Poh *PohFilterer) ParseArbitratorComplete(log types.Log) (*PohArbitratorComplete, error) {
	event := new(PohArbitratorComplete)
	if err := _Poh.contract.UnpackLog(event, "ArbitratorComplete", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PohChallengeResolvedIterator is returned from FilterChallengeResolved and is used to iterate over the raw logs and unpacked data for ChallengeResolved events raised by the Poh contract.
type PohChallengeResolvedIterator struct {
	Event *PohChallengeResolved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PohChallengeResolvedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PohChallengeResolved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PohChallengeResolved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PohChallengeResolvedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PohChallengeResolvedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PohChallengeResolved represents a ChallengeResolved event raised by the Poh contract.
type PohChallengeResolved struct {
	SubmissionID common.Address
	RequestID    *big.Int
	ChallengeID  *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterChallengeResolved is a free log retrieval operation binding the contract event 0xb6759576305cce1591ca803d5fbf22b83b8a7465c093df7b013cb829e98718e1.
//
// Solidity: event ChallengeResolved(address indexed _submissionID, uint256 indexed _requestID, uint256 _challengeID)
func (_Poh *PohFilterer) FilterChallengeResolved(opts *bind.FilterOpts, _submissionID []common.Address, _requestID []*big.Int) (*PohChallengeResolvedIterator, error) {

	var _submissionIDRule []interface{}
	for _, _submissionIDItem := range _submissionID {
		_submissionIDRule = append(_submissionIDRule, _submissionIDItem)
	}
	var _requestIDRule []interface{}
	for _, _requestIDItem := range _requestID {
		_requestIDRule = append(_requestIDRule, _requestIDItem)
	}

	logs, sub, err := _Poh.contract.FilterLogs(opts, "ChallengeResolved", _submissionIDRule, _requestIDRule)
	if err != nil {
		return nil, err
	}
	return &PohChallengeResolvedIterator{contract: _Poh.contract, event: "ChallengeResolved", logs: logs, sub: sub}, nil
}

// WatchChallengeResolved is a free log subscription operation binding the contract event 0xb6759576305cce1591ca803d5fbf22b83b8a7465c093df7b013cb829e98718e1.
//
// Solidity: event ChallengeResolved(address indexed _submissionID, uint256 indexed _requestID, uint256 _challengeID)
func (_Poh *PohFilterer) WatchChallengeResolved(opts *bind.WatchOpts, sink chan<- *PohChallengeResolved, _submissionID []common.Address, _requestID []*big.Int) (event.Subscription, error) {

	var _submissionIDRule []interface{}
	for _, _submissionIDItem := range _submissionID {
		_submissionIDRule = append(_submissionIDRule, _submissionIDItem)
	}
	var _requestIDRule []interface{}
	for _, _requestIDItem := range _requestID {
		_requestIDRule = append(_requestIDRule, _requestIDItem)
	}

	logs, sub, err := _Poh.contract.WatchLogs(opts, "ChallengeResolved", _submissionIDRule, _requestIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PohChallengeResolved)
				if err := _Poh.contract.UnpackLog(event, "ChallengeResolved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseChallengeResolved is a log parse operation binding the contract event 0xb6759576305cce1591ca803d5fbf22b83b8a7465c093df7b013cb829e98718e1.
//
// Solidity: event ChallengeResolved(address indexed _submissionID, uint256 indexed _requestID, uint256 _challengeID)
func (_Poh *PohFilterer) ParseChallengeResolved(log types.Log) (*PohChallengeResolved, error) {
	event := new(PohChallengeResolved)
	if err := _Poh.contract.UnpackLog(event, "ChallengeResolved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PohDisputeIterator is returned from FilterDispute and is used to iterate over the raw logs and unpacked data for Dispute events raised by the Poh contract.
type PohDisputeIterator struct {
	Event *PohDispute // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PohDisputeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PohDispute)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PohDispute)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PohDisputeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PohDisputeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PohDispute represents a Dispute event raised by the Poh contract.
type PohDispute struct {
	Arbitrator      common.Address
	DisputeID       *big.Int
	MetaEvidenceID  *big.Int
	EvidenceGroupID *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterDispute is a free log retrieval operation binding the contract event 0x74baab670a4015ab2f1b467c5252a96141a2573f2908e58a92081e80d3cfde3d.
//
// Solidity: event Dispute(address indexed _arbitrator, uint256 indexed _disputeID, uint256 _metaEvidenceID, uint256 _evidenceGroupID)
func (_Poh *PohFilterer) FilterDispute(opts *bind.FilterOpts, _arbitrator []common.Address, _disputeID []*big.Int) (*PohDisputeIterator, error) {

	var _arbitratorRule []interface{}
	for _, _arbitratorItem := range _arbitrator {
		_arbitratorRule = append(_arbitratorRule, _arbitratorItem)
	}
	var _disputeIDRule []interface{}
	for _, _disputeIDItem := range _disputeID {
		_disputeIDRule = append(_disputeIDRule, _disputeIDItem)
	}

	logs, sub, err := _Poh.contract.FilterLogs(opts, "Dispute", _arbitratorRule, _disputeIDRule)
	if err != nil {
		return nil, err
	}
	return &PohDisputeIterator{contract: _Poh.contract, event: "Dispute", logs: logs, sub: sub}, nil
}

// WatchDispute is a free log subscription operation binding the contract event 0x74baab670a4015ab2f1b467c5252a96141a2573f2908e58a92081e80d3cfde3d.
//
// Solidity: event Dispute(address indexed _arbitrator, uint256 indexed _disputeID, uint256 _metaEvidenceID, uint256 _evidenceGroupID)
func (_Poh *PohFilterer) WatchDispute(opts *bind.WatchOpts, sink chan<- *PohDispute, _arbitrator []common.Address, _disputeID []*big.Int) (event.Subscription, error) {

	var _arbitratorRule []interface{}
	for _, _arbitratorItem := range _arbitrator {
		_arbitratorRule = append(_arbitratorRule, _arbitratorItem)
	}
	var _disputeIDRule []interface{}
	for _, _disputeIDItem := range _disputeID {
		_disputeIDRule = append(_disputeIDRule, _disputeIDItem)
	}

	logs, sub, err := _Poh.contract.WatchLogs(opts, "Dispute", _arbitratorRule, _disputeIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PohDispute)
				if err := _Poh.contract.UnpackLog(event, "Dispute", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDispute is a log parse operation binding the contract event 0x74baab670a4015ab2f1b467c5252a96141a2573f2908e58a92081e80d3cfde3d.
//
// Solidity: event Dispute(address indexed _arbitrator, uint256 indexed _disputeID, uint256 _metaEvidenceID, uint256 _evidenceGroupID)
func (_Poh *PohFilterer) ParseDispute(log types.Log) (*PohDispute, error) {
	event := new(PohDispute)
	if err := _Poh.contract.UnpackLog(event, "Dispute", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PohEvidenceIterator is returned from FilterEvidence and is used to iterate over the raw logs and unpacked data for Evidence events raised by the Poh contract.
type PohEvidenceIterator struct {
	Event *PohEvidence // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PohEvidenceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PohEvidence)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PohEvidence)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PohEvidenceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PohEvidenceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PohEvidence represents a Evidence event raised by the Poh contract.
type PohEvidence struct {
	Arbitrator      common.Address
	EvidenceGroupID *big.Int
	Party           common.Address
	Evidence        string
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterEvidence is a free log retrieval operation binding the contract event 0xdccf2f8b2cc26eafcd61905cba744cff4b81d14740725f6376390dc6298a6a3c.
//
// Solidity: event Evidence(address indexed _arbitrator, uint256 indexed _evidenceGroupID, address indexed _party, string _evidence)
func (_Poh *PohFilterer) FilterEvidence(opts *bind.FilterOpts, _arbitrator []common.Address, _evidenceGroupID []*big.Int, _party []common.Address) (*PohEvidenceIterator, error) {

	var _arbitratorRule []interface{}
	for _, _arbitratorItem := range _arbitrator {
		_arbitratorRule = append(_arbitratorRule, _arbitratorItem)
	}
	var _evidenceGroupIDRule []interface{}
	for _, _evidenceGroupIDItem := range _evidenceGroupID {
		_evidenceGroupIDRule = append(_evidenceGroupIDRule, _evidenceGroupIDItem)
	}
	var _partyRule []interface{}
	for _, _partyItem := range _party {
		_partyRule = append(_partyRule, _partyItem)
	}

	logs, sub, err := _Poh.contract.FilterLogs(opts, "Evidence", _arbitratorRule, _evidenceGroupIDRule, _partyRule)
	if err != nil {
		return nil, err
	}
	return &PohEvidenceIterator{contract: _Poh.contract, event: "Evidence", logs: logs, sub: sub}, nil
}

// WatchEvidence is a free log subscription operation binding the contract event 0xdccf2f8b2cc26eafcd61905cba744cff4b81d14740725f6376390dc6298a6a3c.
//
// Solidity: event Evidence(address indexed _arbitrator, uint256 indexed _evidenceGroupID, address indexed _party, string _evidence)
func (_Poh *PohFilterer) WatchEvidence(opts *bind.WatchOpts, sink chan<- *PohEvidence, _arbitrator []common.Address, _evidenceGroupID []*big.Int, _party []common.Address) (event.Subscription, error) {

	var _arbitratorRule []interface{}
	for _, _arbitratorItem := range _arbitrator {
		_arbitratorRule = append(_arbitratorRule, _arbitratorItem)
	}
	var _evidenceGroupIDRule []interface{}
	for _, _evidenceGroupIDItem := range _evidenceGroupID {
		_evidenceGroupIDRule = append(_evidenceGroupIDRule, _evidenceGroupIDItem)
	}
	var _partyRule []interface{}
	for _, _partyItem := range _party {
		_partyRule = append(_partyRule, _partyItem)
	}

	logs, sub, err := _Poh.contract.WatchLogs(opts, "Evidence", _arbitratorRule, _evidenceGroupIDRule, _partyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PohEvidence)
				if err := _Poh.contract.UnpackLog(event, "Evidence", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseEvidence is a log parse operation binding the contract event 0xdccf2f8b2cc26eafcd61905cba744cff4b81d14740725f6376390dc6298a6a3c.
//
// Solidity: event Evidence(address indexed _arbitrator, uint256 indexed _evidenceGroupID, address indexed _party, string _evidence)
func (_Poh *PohFilterer) ParseEvidence(log types.Log) (*PohEvidence, error) {
	event := new(PohEvidence)
	if err := _Poh.contract.UnpackLog(event, "Evidence", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PohHasPaidAppealFeeIterator is returned from FilterHasPaidAppealFee and is used to iterate over the raw logs and unpacked data for HasPaidAppealFee events raised by the Poh contract.
type PohHasPaidAppealFeeIterator struct {
	Event *PohHasPaidAppealFee // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PohHasPaidAppealFeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PohHasPaidAppealFee)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PohHasPaidAppealFee)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PohHasPaidAppealFeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PohHasPaidAppealFeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PohHasPaidAppealFee represents a HasPaidAppealFee event raised by the Poh contract.
type PohHasPaidAppealFee struct {
	SubmissionID common.Address
	ChallengeID  *big.Int
	Side         uint8
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterHasPaidAppealFee is a free log retrieval operation binding the contract event 0x642c5385391e218917cf658cc1365fdffc24183646b08706ed51281972d1dc80.
//
// Solidity: event HasPaidAppealFee(address indexed _submissionID, uint256 indexed _challengeID, uint8 _side)
func (_Poh *PohFilterer) FilterHasPaidAppealFee(opts *bind.FilterOpts, _submissionID []common.Address, _challengeID []*big.Int) (*PohHasPaidAppealFeeIterator, error) {

	var _submissionIDRule []interface{}
	for _, _submissionIDItem := range _submissionID {
		_submissionIDRule = append(_submissionIDRule, _submissionIDItem)
	}
	var _challengeIDRule []interface{}
	for _, _challengeIDItem := range _challengeID {
		_challengeIDRule = append(_challengeIDRule, _challengeIDItem)
	}

	logs, sub, err := _Poh.contract.FilterLogs(opts, "HasPaidAppealFee", _submissionIDRule, _challengeIDRule)
	if err != nil {
		return nil, err
	}
	return &PohHasPaidAppealFeeIterator{contract: _Poh.contract, event: "HasPaidAppealFee", logs: logs, sub: sub}, nil
}

// WatchHasPaidAppealFee is a free log subscription operation binding the contract event 0x642c5385391e218917cf658cc1365fdffc24183646b08706ed51281972d1dc80.
//
// Solidity: event HasPaidAppealFee(address indexed _submissionID, uint256 indexed _challengeID, uint8 _side)
func (_Poh *PohFilterer) WatchHasPaidAppealFee(opts *bind.WatchOpts, sink chan<- *PohHasPaidAppealFee, _submissionID []common.Address, _challengeID []*big.Int) (event.Subscription, error) {

	var _submissionIDRule []interface{}
	for _, _submissionIDItem := range _submissionID {
		_submissionIDRule = append(_submissionIDRule, _submissionIDItem)
	}
	var _challengeIDRule []interface{}
	for _, _challengeIDItem := range _challengeID {
		_challengeIDRule = append(_challengeIDRule, _challengeIDItem)
	}

	logs, sub, err := _Poh.contract.WatchLogs(opts, "HasPaidAppealFee", _submissionIDRule, _challengeIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PohHasPaidAppealFee)
				if err := _Poh.contract.UnpackLog(event, "HasPaidAppealFee", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseHasPaidAppealFee is a log parse operation binding the contract event 0x642c5385391e218917cf658cc1365fdffc24183646b08706ed51281972d1dc80.
//
// Solidity: event HasPaidAppealFee(address indexed _submissionID, uint256 indexed _challengeID, uint8 _side)
func (_Poh *PohFilterer) ParseHasPaidAppealFee(log types.Log) (*PohHasPaidAppealFee, error) {
	event := new(PohHasPaidAppealFee)
	if err := _Poh.contract.UnpackLog(event, "HasPaidAppealFee", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PohMetaEvidenceIterator is returned from FilterMetaEvidence and is used to iterate over the raw logs and unpacked data for MetaEvidence events raised by the Poh contract.
type PohMetaEvidenceIterator struct {
	Event *PohMetaEvidence // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PohMetaEvidenceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PohMetaEvidence)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PohMetaEvidence)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PohMetaEvidenceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PohMetaEvidenceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PohMetaEvidence represents a MetaEvidence event raised by the Poh contract.
type PohMetaEvidence struct {
	MetaEvidenceID *big.Int
	Evidence       string
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterMetaEvidence is a free log retrieval operation binding the contract event 0x61606860eb6c87306811e2695215385101daab53bd6ab4e9f9049aead9363c7d.
//
// Solidity: event MetaEvidence(uint256 indexed _metaEvidenceID, string _evidence)
func (_Poh *PohFilterer) FilterMetaEvidence(opts *bind.FilterOpts, _metaEvidenceID []*big.Int) (*PohMetaEvidenceIterator, error) {

	var _metaEvidenceIDRule []interface{}
	for _, _metaEvidenceIDItem := range _metaEvidenceID {
		_metaEvidenceIDRule = append(_metaEvidenceIDRule, _metaEvidenceIDItem)
	}

	logs, sub, err := _Poh.contract.FilterLogs(opts, "MetaEvidence", _metaEvidenceIDRule)
	if err != nil {
		return nil, err
	}
	return &PohMetaEvidenceIterator{contract: _Poh.contract, event: "MetaEvidence", logs: logs, sub: sub}, nil
}

// WatchMetaEvidence is a free log subscription operation binding the contract event 0x61606860eb6c87306811e2695215385101daab53bd6ab4e9f9049aead9363c7d.
//
// Solidity: event MetaEvidence(uint256 indexed _metaEvidenceID, string _evidence)
func (_Poh *PohFilterer) WatchMetaEvidence(opts *bind.WatchOpts, sink chan<- *PohMetaEvidence, _metaEvidenceID []*big.Int) (event.Subscription, error) {

	var _metaEvidenceIDRule []interface{}
	for _, _metaEvidenceIDItem := range _metaEvidenceID {
		_metaEvidenceIDRule = append(_metaEvidenceIDRule, _metaEvidenceIDItem)
	}

	logs, sub, err := _Poh.contract.WatchLogs(opts, "MetaEvidence", _metaEvidenceIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PohMetaEvidence)
				if err := _Poh.contract.UnpackLog(event, "MetaEvidence", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMetaEvidence is a log parse operation binding the contract event 0x61606860eb6c87306811e2695215385101daab53bd6ab4e9f9049aead9363c7d.
//
// Solidity: event MetaEvidence(uint256 indexed _metaEvidenceID, string _evidence)
func (_Poh *PohFilterer) ParseMetaEvidence(log types.Log) (*PohMetaEvidence, error) {
	event := new(PohMetaEvidence)
	if err := _Poh.contract.UnpackLog(event, "MetaEvidence", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PohReapplySubmissionIterator is returned from FilterReapplySubmission and is used to iterate over the raw logs and unpacked data for ReapplySubmission events raised by the Poh contract.
type PohReapplySubmissionIterator struct {
	Event *PohReapplySubmission // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PohReapplySubmissionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PohReapplySubmission)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PohReapplySubmission)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PohReapplySubmissionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PohReapplySubmissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PohReapplySubmission represents a ReapplySubmission event raised by the Poh contract.
type PohReapplySubmission struct {
	SubmissionID common.Address
	RequestID    *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterReapplySubmission is a free log retrieval operation binding the contract event 0xf6cfccc832db8edf362f395f01d696c7da2db840d97eb1ec1ac44c980143990e.
//
// Solidity: event ReapplySubmission(address indexed _submissionID, uint256 _requestID)
func (_Poh *PohFilterer) FilterReapplySubmission(opts *bind.FilterOpts, _submissionID []common.Address) (*PohReapplySubmissionIterator, error) {

	var _submissionIDRule []interface{}
	for _, _submissionIDItem := range _submissionID {
		_submissionIDRule = append(_submissionIDRule, _submissionIDItem)
	}

	logs, sub, err := _Poh.contract.FilterLogs(opts, "ReapplySubmission", _submissionIDRule)
	if err != nil {
		return nil, err
	}
	return &PohReapplySubmissionIterator{contract: _Poh.contract, event: "ReapplySubmission", logs: logs, sub: sub}, nil
}

// WatchReapplySubmission is a free log subscription operation binding the contract event 0xf6cfccc832db8edf362f395f01d696c7da2db840d97eb1ec1ac44c980143990e.
//
// Solidity: event ReapplySubmission(address indexed _submissionID, uint256 _requestID)
func (_Poh *PohFilterer) WatchReapplySubmission(opts *bind.WatchOpts, sink chan<- *PohReapplySubmission, _submissionID []common.Address) (event.Subscription, error) {

	var _submissionIDRule []interface{}
	for _, _submissionIDItem := range _submissionID {
		_submissionIDRule = append(_submissionIDRule, _submissionIDItem)
	}

	logs, sub, err := _Poh.contract.WatchLogs(opts, "ReapplySubmission", _submissionIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PohReapplySubmission)
				if err := _Poh.contract.UnpackLog(event, "ReapplySubmission", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseReapplySubmission is a log parse operation binding the contract event 0xf6cfccc832db8edf362f395f01d696c7da2db840d97eb1ec1ac44c980143990e.
//
// Solidity: event ReapplySubmission(address indexed _submissionID, uint256 _requestID)
func (_Poh *PohFilterer) ParseReapplySubmission(log types.Log) (*PohReapplySubmission, error) {
	event := new(PohReapplySubmission)
	if err := _Poh.contract.UnpackLog(event, "ReapplySubmission", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PohRemoveSubmissionIterator is returned from FilterRemoveSubmission and is used to iterate over the raw logs and unpacked data for RemoveSubmission events raised by the Poh contract.
type PohRemoveSubmissionIterator struct {
	Event *PohRemoveSubmission // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PohRemoveSubmissionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PohRemoveSubmission)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PohRemoveSubmission)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PohRemoveSubmissionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PohRemoveSubmissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PohRemoveSubmission represents a RemoveSubmission event raised by the Poh contract.
type PohRemoveSubmission struct {
	Requester    common.Address
	SubmissionID common.Address
	RequestID    *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRemoveSubmission is a free log retrieval operation binding the contract event 0xd63ca5272e9e07a30a33cc438f956428bf02359db9798ce3fe30140dadf8d741.
//
// Solidity: event RemoveSubmission(address indexed _requester, address indexed _submissionID, uint256 _requestID)
func (_Poh *PohFilterer) FilterRemoveSubmission(opts *bind.FilterOpts, _requester []common.Address, _submissionID []common.Address) (*PohRemoveSubmissionIterator, error) {

	var _requesterRule []interface{}
	for _, _requesterItem := range _requester {
		_requesterRule = append(_requesterRule, _requesterItem)
	}
	var _submissionIDRule []interface{}
	for _, _submissionIDItem := range _submissionID {
		_submissionIDRule = append(_submissionIDRule, _submissionIDItem)
	}

	logs, sub, err := _Poh.contract.FilterLogs(opts, "RemoveSubmission", _requesterRule, _submissionIDRule)
	if err != nil {
		return nil, err
	}
	return &PohRemoveSubmissionIterator{contract: _Poh.contract, event: "RemoveSubmission", logs: logs, sub: sub}, nil
}

// WatchRemoveSubmission is a free log subscription operation binding the contract event 0xd63ca5272e9e07a30a33cc438f956428bf02359db9798ce3fe30140dadf8d741.
//
// Solidity: event RemoveSubmission(address indexed _requester, address indexed _submissionID, uint256 _requestID)
func (_Poh *PohFilterer) WatchRemoveSubmission(opts *bind.WatchOpts, sink chan<- *PohRemoveSubmission, _requester []common.Address, _submissionID []common.Address) (event.Subscription, error) {

	var _requesterRule []interface{}
	for _, _requesterItem := range _requester {
		_requesterRule = append(_requesterRule, _requesterItem)
	}
	var _submissionIDRule []interface{}
	for _, _submissionIDItem := range _submissionID {
		_submissionIDRule = append(_submissionIDRule, _submissionIDItem)
	}

	logs, sub, err := _Poh.contract.WatchLogs(opts, "RemoveSubmission", _requesterRule, _submissionIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PohRemoveSubmission)
				if err := _Poh.contract.UnpackLog(event, "RemoveSubmission", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRemoveSubmission is a log parse operation binding the contract event 0xd63ca5272e9e07a30a33cc438f956428bf02359db9798ce3fe30140dadf8d741.
//
// Solidity: event RemoveSubmission(address indexed _requester, address indexed _submissionID, uint256 _requestID)
func (_Poh *PohFilterer) ParseRemoveSubmission(log types.Log) (*PohRemoveSubmission, error) {
	event := new(PohRemoveSubmission)
	if err := _Poh.contract.UnpackLog(event, "RemoveSubmission", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PohRulingIterator is returned from FilterRuling and is used to iterate over the raw logs and unpacked data for Ruling events raised by the Poh contract.
type PohRulingIterator struct {
	Event *PohRuling // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PohRulingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PohRuling)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PohRuling)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PohRulingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PohRulingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PohRuling represents a Ruling event raised by the Poh contract.
type PohRuling struct {
	Arbitrator common.Address
	DisputeID  *big.Int
	Ruling     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRuling is a free log retrieval operation binding the contract event 0x394027a5fa6e098a1191094d1719d6929b9abc535fcc0c8f448d6a4e75622276.
//
// Solidity: event Ruling(address indexed _arbitrator, uint256 indexed _disputeID, uint256 _ruling)
func (_Poh *PohFilterer) FilterRuling(opts *bind.FilterOpts, _arbitrator []common.Address, _disputeID []*big.Int) (*PohRulingIterator, error) {

	var _arbitratorRule []interface{}
	for _, _arbitratorItem := range _arbitrator {
		_arbitratorRule = append(_arbitratorRule, _arbitratorItem)
	}
	var _disputeIDRule []interface{}
	for _, _disputeIDItem := range _disputeID {
		_disputeIDRule = append(_disputeIDRule, _disputeIDItem)
	}

	logs, sub, err := _Poh.contract.FilterLogs(opts, "Ruling", _arbitratorRule, _disputeIDRule)
	if err != nil {
		return nil, err
	}
	return &PohRulingIterator{contract: _Poh.contract, event: "Ruling", logs: logs, sub: sub}, nil
}

// WatchRuling is a free log subscription operation binding the contract event 0x394027a5fa6e098a1191094d1719d6929b9abc535fcc0c8f448d6a4e75622276.
//
// Solidity: event Ruling(address indexed _arbitrator, uint256 indexed _disputeID, uint256 _ruling)
func (_Poh *PohFilterer) WatchRuling(opts *bind.WatchOpts, sink chan<- *PohRuling, _arbitrator []common.Address, _disputeID []*big.Int) (event.Subscription, error) {

	var _arbitratorRule []interface{}
	for _, _arbitratorItem := range _arbitrator {
		_arbitratorRule = append(_arbitratorRule, _arbitratorItem)
	}
	var _disputeIDRule []interface{}
	for _, _disputeIDItem := range _disputeID {
		_disputeIDRule = append(_disputeIDRule, _disputeIDItem)
	}

	logs, sub, err := _Poh.contract.WatchLogs(opts, "Ruling", _arbitratorRule, _disputeIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PohRuling)
				if err := _Poh.contract.UnpackLog(event, "Ruling", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRuling is a log parse operation binding the contract event 0x394027a5fa6e098a1191094d1719d6929b9abc535fcc0c8f448d6a4e75622276.
//
// Solidity: event Ruling(address indexed _arbitrator, uint256 indexed _disputeID, uint256 _ruling)
func (_Poh *PohFilterer) ParseRuling(log types.Log) (*PohRuling, error) {
	event := new(PohRuling)
	if err := _Poh.contract.UnpackLog(event, "Ruling", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PohSubmissionChallengedIterator is returned from FilterSubmissionChallenged and is used to iterate over the raw logs and unpacked data for SubmissionChallenged events raised by the Poh contract.
type PohSubmissionChallengedIterator struct {
	Event *PohSubmissionChallenged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PohSubmissionChallengedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PohSubmissionChallenged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PohSubmissionChallenged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PohSubmissionChallengedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PohSubmissionChallengedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PohSubmissionChallenged represents a SubmissionChallenged event raised by the Poh contract.
type PohSubmissionChallenged struct {
	SubmissionID common.Address
	RequestID    *big.Int
	ChallengeID  *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterSubmissionChallenged is a free log retrieval operation binding the contract event 0x28ec07f413c7805003c29837f7f1b3799f30f7f87a1e1b5b9aaee94f665218ac.
//
// Solidity: event SubmissionChallenged(address indexed _submissionID, uint256 indexed _requestID, uint256 _challengeID)
func (_Poh *PohFilterer) FilterSubmissionChallenged(opts *bind.FilterOpts, _submissionID []common.Address, _requestID []*big.Int) (*PohSubmissionChallengedIterator, error) {

	var _submissionIDRule []interface{}
	for _, _submissionIDItem := range _submissionID {
		_submissionIDRule = append(_submissionIDRule, _submissionIDItem)
	}
	var _requestIDRule []interface{}
	for _, _requestIDItem := range _requestID {
		_requestIDRule = append(_requestIDRule, _requestIDItem)
	}

	logs, sub, err := _Poh.contract.FilterLogs(opts, "SubmissionChallenged", _submissionIDRule, _requestIDRule)
	if err != nil {
		return nil, err
	}
	return &PohSubmissionChallengedIterator{contract: _Poh.contract, event: "SubmissionChallenged", logs: logs, sub: sub}, nil
}

// WatchSubmissionChallenged is a free log subscription operation binding the contract event 0x28ec07f413c7805003c29837f7f1b3799f30f7f87a1e1b5b9aaee94f665218ac.
//
// Solidity: event SubmissionChallenged(address indexed _submissionID, uint256 indexed _requestID, uint256 _challengeID)
func (_Poh *PohFilterer) WatchSubmissionChallenged(opts *bind.WatchOpts, sink chan<- *PohSubmissionChallenged, _submissionID []common.Address, _requestID []*big.Int) (event.Subscription, error) {

	var _submissionIDRule []interface{}
	for _, _submissionIDItem := range _submissionID {
		_submissionIDRule = append(_submissionIDRule, _submissionIDItem)
	}
	var _requestIDRule []interface{}
	for _, _requestIDItem := range _requestID {
		_requestIDRule = append(_requestIDRule, _requestIDItem)
	}

	logs, sub, err := _Poh.contract.WatchLogs(opts, "SubmissionChallenged", _submissionIDRule, _requestIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PohSubmissionChallenged)
				if err := _Poh.contract.UnpackLog(event, "SubmissionChallenged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSubmissionChallenged is a log parse operation binding the contract event 0x28ec07f413c7805003c29837f7f1b3799f30f7f87a1e1b5b9aaee94f665218ac.
//
// Solidity: event SubmissionChallenged(address indexed _submissionID, uint256 indexed _requestID, uint256 _challengeID)
func (_Poh *PohFilterer) ParseSubmissionChallenged(log types.Log) (*PohSubmissionChallenged, error) {
	event := new(PohSubmissionChallenged)
	if err := _Poh.contract.UnpackLog(event, "SubmissionChallenged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PohVouchAddedIterator is returned from FilterVouchAdded and is used to iterate over the raw logs and unpacked data for VouchAdded events raised by the Poh contract.
type PohVouchAddedIterator struct {
	Event *PohVouchAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PohVouchAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PohVouchAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PohVouchAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PohVouchAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PohVouchAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PohVouchAdded represents a VouchAdded event raised by the Poh contract.
type PohVouchAdded struct {
	SubmissionID common.Address
	Voucher      common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterVouchAdded is a free log retrieval operation binding the contract event 0xc5482a4357210d662eeea6fa4589e76e47bdb84517a3b9da4f8b7576ae001701.
//
// Solidity: event VouchAdded(address indexed _submissionID, address indexed _voucher)
func (_Poh *PohFilterer) FilterVouchAdded(opts *bind.FilterOpts, _submissionID []common.Address, _voucher []common.Address) (*PohVouchAddedIterator, error) {

	var _submissionIDRule []interface{}
	for _, _submissionIDItem := range _submissionID {
		_submissionIDRule = append(_submissionIDRule, _submissionIDItem)
	}
	var _voucherRule []interface{}
	for _, _voucherItem := range _voucher {
		_voucherRule = append(_voucherRule, _voucherItem)
	}

	logs, sub, err := _Poh.contract.FilterLogs(opts, "VouchAdded", _submissionIDRule, _voucherRule)
	if err != nil {
		return nil, err
	}
	return &PohVouchAddedIterator{contract: _Poh.contract, event: "VouchAdded", logs: logs, sub: sub}, nil
}

// WatchVouchAdded is a free log subscription operation binding the contract event 0xc5482a4357210d662eeea6fa4589e76e47bdb84517a3b9da4f8b7576ae001701.
//
// Solidity: event VouchAdded(address indexed _submissionID, address indexed _voucher)
func (_Poh *PohFilterer) WatchVouchAdded(opts *bind.WatchOpts, sink chan<- *PohVouchAdded, _submissionID []common.Address, _voucher []common.Address) (event.Subscription, error) {

	var _submissionIDRule []interface{}
	for _, _submissionIDItem := range _submissionID {
		_submissionIDRule = append(_submissionIDRule, _submissionIDItem)
	}
	var _voucherRule []interface{}
	for _, _voucherItem := range _voucher {
		_voucherRule = append(_voucherRule, _voucherItem)
	}

	logs, sub, err := _Poh.contract.WatchLogs(opts, "VouchAdded", _submissionIDRule, _voucherRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PohVouchAdded)
				if err := _Poh.contract.UnpackLog(event, "VouchAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseVouchAdded is a log parse operation binding the contract event 0xc5482a4357210d662eeea6fa4589e76e47bdb84517a3b9da4f8b7576ae001701.
//
// Solidity: event VouchAdded(address indexed _submissionID, address indexed _voucher)
func (_Poh *PohFilterer) ParseVouchAdded(log types.Log) (*PohVouchAdded, error) {
	event := new(PohVouchAdded)
	if err := _Poh.contract.UnpackLog(event, "VouchAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PohVouchRemovedIterator is returned from FilterVouchRemoved and is used to iterate over the raw logs and unpacked data for VouchRemoved events raised by the Poh contract.
type PohVouchRemovedIterator struct {
	Event *PohVouchRemoved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PohVouchRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PohVouchRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PohVouchRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PohVouchRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PohVouchRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PohVouchRemoved represents a VouchRemoved event raised by the Poh contract.
type PohVouchRemoved struct {
	SubmissionID common.Address
	Voucher      common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterVouchRemoved is a free log retrieval operation binding the contract event 0xd6f782ad61ba30c1b39b77f1bf37061e68733360288a441e70d08c0a91b0cbc5.
//
// Solidity: event VouchRemoved(address indexed _submissionID, address indexed _voucher)
func (_Poh *PohFilterer) FilterVouchRemoved(opts *bind.FilterOpts, _submissionID []common.Address, _voucher []common.Address) (*PohVouchRemovedIterator, error) {

	var _submissionIDRule []interface{}
	for _, _submissionIDItem := range _submissionID {
		_submissionIDRule = append(_submissionIDRule, _submissionIDItem)
	}
	var _voucherRule []interface{}
	for _, _voucherItem := range _voucher {
		_voucherRule = append(_voucherRule, _voucherItem)
	}

	logs, sub, err := _Poh.contract.FilterLogs(opts, "VouchRemoved", _submissionIDRule, _voucherRule)
	if err != nil {
		return nil, err
	}
	return &PohVouchRemovedIterator{contract: _Poh.contract, event: "VouchRemoved", logs: logs, sub: sub}, nil
}

// WatchVouchRemoved is a free log subscription operation binding the contract event 0xd6f782ad61ba30c1b39b77f1bf37061e68733360288a441e70d08c0a91b0cbc5.
//
// Solidity: event VouchRemoved(address indexed _submissionID, address indexed _voucher)
func (_Poh *PohFilterer) WatchVouchRemoved(opts *bind.WatchOpts, sink chan<- *PohVouchRemoved, _submissionID []common.Address, _voucher []common.Address) (event.Subscription, error) {

	var _submissionIDRule []interface{}
	for _, _submissionIDItem := range _submissionID {
		_submissionIDRule = append(_submissionIDRule, _submissionIDItem)
	}
	var _voucherRule []interface{}
	for _, _voucherItem := range _voucher {
		_voucherRule = append(_voucherRule, _voucherItem)
	}

	logs, sub, err := _Poh.contract.WatchLogs(opts, "VouchRemoved", _submissionIDRule, _voucherRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PohVouchRemoved)
				if err := _Poh.contract.UnpackLog(event, "VouchRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseVouchRemoved is a log parse operation binding the contract event 0xd6f782ad61ba30c1b39b77f1bf37061e68733360288a441e70d08c0a91b0cbc5.
//
// Solidity: event VouchRemoved(address indexed _submissionID, address indexed _voucher)
func (_Poh *PohFilterer) ParseVouchRemoved(log types.Log) (*PohVouchRemoved, error) {
	event := new(PohVouchRemoved)
	if err := _Poh.contract.UnpackLog(event, "VouchRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
