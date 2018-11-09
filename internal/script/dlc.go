package script

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// ContractExecutionDelay is a delay used in ContractExecutionScript
const ContractExecutionDelay = 144

// ContractExecutionScript returns a contract execution script.
//
// Script Code:
//  OP_IF
//    <public key a + message public key>
//  OP_ELSE
//    delay(fix 144)
//    OP_CHECKSEQUENCEVERIFY
//    OP_DROP
//    <public key b>
//  OP_ENDIF
//  OP_CHECKSIG
//
// The if block can be passed when the contractor A has a valid oracle's sign to the message.
// But if the contractor sends this transaction without the oracle's valid sign,
// the else block will be used by the other party B after the delay time (1 day approximately).
// Please check the original paper for more details.
//
// https://adiabat.github.io/dlc.pdf
func ContractExecutionScript(puba, pubb, pubm *btcec.PublicKey) ([]byte, error) {
	// pub key a + message pub key
	pubam := &btcec.PublicKey{}
	pubam.X, pubam.Y = btcec.S256().Add(puba.X, puba.Y, pubm.X, pubm.Y)

	delay := uint16(ContractExecutionDelay)
	csvflg := uint32(0x00000000)
	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_IF)
	builder.AddData(pubam.SerializeCompressed())
	builder.AddOp(txscript.OP_ELSE)
	builder.AddInt64(int64(delay) + int64(csvflg))
	builder.AddOp(txscript.OP_CHECKSEQUENCEVERIFY)
	builder.AddOp(txscript.OP_DROP)
	builder.AddData(pubb.SerializeCompressed())
	builder.AddOp(txscript.OP_ENDIF)
	builder.AddOp(txscript.OP_CHECKSIG)
	return builder.Script()
}

// WitnessForCEScript constructs a witness that unlocks a contract execution script.
// This function use the OP_IF block
func WitnessForCEScript(
	sign []byte, script []byte) wire.TxWitness {
	return wire.TxWitness{sign, []byte{1}, script}
}

// WitnessForCEScriptAfterDelay constructs a witness that unlocks a contract execution script.
// This function use the OP_ELSE block that can be valid after the delay
func WitnessForCEScriptAfterDelay(
	sign []byte, script []byte) wire.TxWitness {
	return wire.TxWitness{sign, []byte{}, script}
}
