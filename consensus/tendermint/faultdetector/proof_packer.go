package faultdetector


/*
	Pack
	type Proof struct {
		Rule     Rule
		Message  message
		Evidence []message
	}

	into a byte array which is represented by above elements for static unpacking.

	type ParsableProof struct {
		Rule uint8;				 // take the 1st byte as rule id.
		NumberOfMsg uint32;		 // take the next 4 bytes as number of msgs.
		LengthOfEach []uint32;   // take the next "NumberOfMsg" * 4 bytes as the set of lengths for each message.
		MessagePayloads []bytes; // take the left bytes as payload of each message tightly packed, slotted by own lengths.
	}

 */
func packProof(proof *Proof) []byte {
	return nil
}

func unpackProof(packedProof []byte) *Proof {
	return nil
}


