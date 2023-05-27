// Code generated by github.com/whyrusleeping/cbor-gen. DO NOT EDIT.

package assets

import (
	"fmt"
	"io"
	"math"
	"sort"

	cid "github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
	xerrors "golang.org/x/xerrors"
)

var _ = xerrors.Errorf
var _ = cid.Undef
var _ = math.E
var _ = sort.Sort

func (t *AssetPullingInfo) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write([]byte{176}); err != nil {
		return err
	}

	// t.CID (string) (string)
	if len("CID") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"CID\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("CID"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("CID")); err != nil {
		return err
	}

	if len(t.CID) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.CID was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.CID))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.CID)); err != nil {
		return err
	}

	// t.Hash (assets.AssetHash) (string)
	if len("Hash") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Hash\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("Hash"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Hash")); err != nil {
		return err
	}

	if len(t.Hash) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Hash was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.Hash))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Hash)); err != nil {
		return err
	}

	// t.Size (int64) (int64)
	if len("Size") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Size\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("Size"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Size")); err != nil {
		return err
	}

	if t.Size >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Size)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.Size-1)); err != nil {
			return err
		}
	}

	// t.State (assets.AssetState) (string)
	if len("State") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"State\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("State"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("State")); err != nil {
		return err
	}

	if len(t.State) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.State was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.State))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.State)); err != nil {
		return err
	}

	// t.Blocks (int64) (int64)
	if len("Blocks") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Blocks\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("Blocks"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Blocks")); err != nil {
		return err
	}

	if t.Blocks >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Blocks)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.Blocks-1)); err != nil {
			return err
		}
	}

	// t.Details (string) (string)
	if len("Details") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Details\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("Details"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Details")); err != nil {
		return err
	}

	if len(t.Details) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Details was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.Details))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Details)); err != nil {
		return err
	}

	// t.Bandwidth (int64) (int64)
	if len("Bandwidth") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Bandwidth\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("Bandwidth"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Bandwidth")); err != nil {
		return err
	}

	if t.Bandwidth >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Bandwidth)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.Bandwidth-1)); err != nil {
			return err
		}
	}

	// t.Requester (string) (string)
	if len("Requester") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Requester\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("Requester"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Requester")); err != nil {
		return err
	}

	if len(t.Requester) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Requester was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.Requester))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Requester)); err != nil {
		return err
	}

	// t.RetryCount (int64) (int64)
	if len("RetryCount") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"RetryCount\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("RetryCount"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("RetryCount")); err != nil {
		return err
	}

	if t.RetryCount >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.RetryCount)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.RetryCount-1)); err != nil {
			return err
		}
	}

	// t.EdgeReplicas (int64) (int64)
	if len("EdgeReplicas") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"EdgeReplicas\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("EdgeReplicas"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("EdgeReplicas")); err != nil {
		return err
	}

	if t.EdgeReplicas >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.EdgeReplicas)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.EdgeReplicas-1)); err != nil {
			return err
		}
	}

	// t.EdgeWaitings (int64) (int64)
	if len("EdgeWaitings") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"EdgeWaitings\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("EdgeWaitings"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("EdgeWaitings")); err != nil {
		return err
	}

	if t.EdgeWaitings >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.EdgeWaitings)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.EdgeWaitings-1)); err != nil {
			return err
		}
	}

	// t.CandidateReplicas (int64) (int64)
	if len("CandidateReplicas") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"CandidateReplicas\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("CandidateReplicas"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("CandidateReplicas")); err != nil {
		return err
	}

	if t.CandidateReplicas >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.CandidateReplicas)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.CandidateReplicas-1)); err != nil {
			return err
		}
	}

	// t.CandidateWaitings (int64) (int64)
	if len("CandidateWaitings") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"CandidateWaitings\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("CandidateWaitings"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("CandidateWaitings")); err != nil {
		return err
	}

	if t.CandidateWaitings >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.CandidateWaitings)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.CandidateWaitings-1)); err != nil {
			return err
		}
	}

	// t.ReplenishReplicas (int64) (int64)
	if len("ReplenishReplicas") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"ReplenishReplicas\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("ReplenishReplicas"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("ReplenishReplicas")); err != nil {
		return err
	}

	if t.ReplenishReplicas >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.ReplenishReplicas)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.ReplenishReplicas-1)); err != nil {
			return err
		}
	}

	// t.EdgeReplicaSucceeds ([]string) (slice)
	if len("EdgeReplicaSucceeds") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"EdgeReplicaSucceeds\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("EdgeReplicaSucceeds"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("EdgeReplicaSucceeds")); err != nil {
		return err
	}

	if len(t.EdgeReplicaSucceeds) > cbg.MaxLength {
		return xerrors.Errorf("Slice value in field t.EdgeReplicaSucceeds was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.EdgeReplicaSucceeds))); err != nil {
		return err
	}
	for _, v := range t.EdgeReplicaSucceeds {
		if len(v) > cbg.MaxLength {
			return xerrors.Errorf("Value in field v was too long")
		}

		if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(v))); err != nil {
			return err
		}
		if _, err := io.WriteString(w, string(v)); err != nil {
			return err
		}
	}

	// t.CandidateReplicaSucceeds ([]string) (slice)
	if len("CandidateReplicaSucceeds") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"CandidateReplicaSucceeds\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("CandidateReplicaSucceeds"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("CandidateReplicaSucceeds")); err != nil {
		return err
	}

	if len(t.CandidateReplicaSucceeds) > cbg.MaxLength {
		return xerrors.Errorf("Slice value in field t.CandidateReplicaSucceeds was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.CandidateReplicaSucceeds))); err != nil {
		return err
	}
	for _, v := range t.CandidateReplicaSucceeds {
		if len(v) > cbg.MaxLength {
			return xerrors.Errorf("Value in field v was too long")
		}

		if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(v))); err != nil {
			return err
		}
		if _, err := io.WriteString(w, string(v)); err != nil {
			return err
		}
	}
	return nil
}

func (t *AssetPullingInfo) UnmarshalCBOR(r io.Reader) (err error) {
	*t = AssetPullingInfo{}

	cr := cbg.NewCborReader(r)

	maj, extra, err := cr.ReadHeader()
	if err != nil {
		return err
	}
	defer func() {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	if maj != cbg.MajMap {
		return fmt.Errorf("cbor input should be of type map")
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("AssetPullingInfo: map struct too large (%d)", extra)
	}

	var name string
	n := extra

	for i := uint64(0); i < n; i++ {

		{
			sval, err := cbg.ReadString(cr)
			if err != nil {
				return err
			}

			name = string(sval)
		}

		switch name {
		// t.CID (string) (string)
		case "CID":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.CID = string(sval)
			}
			// t.Hash (assets.AssetHash) (string)
		case "Hash":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.Hash = AssetHash(sval)
			}
			// t.Size (int64) (int64)
		case "Size":
			{
				maj, extra, err := cr.ReadHeader()
				var extraI int64
				if err != nil {
					return err
				}
				switch maj {
				case cbg.MajUnsignedInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 positive overflow")
					}
				case cbg.MajNegativeInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 negative overflow")
					}
					extraI = -1 - extraI
				default:
					return fmt.Errorf("wrong type for int64 field: %d", maj)
				}

				t.Size = int64(extraI)
			}
			// t.State (assets.AssetState) (string)
		case "State":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.State = AssetState(sval)
			}
			// t.Blocks (int64) (int64)
		case "Blocks":
			{
				maj, extra, err := cr.ReadHeader()
				var extraI int64
				if err != nil {
					return err
				}
				switch maj {
				case cbg.MajUnsignedInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 positive overflow")
					}
				case cbg.MajNegativeInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 negative overflow")
					}
					extraI = -1 - extraI
				default:
					return fmt.Errorf("wrong type for int64 field: %d", maj)
				}

				t.Blocks = int64(extraI)
			}
			// t.Details (string) (string)
		case "Details":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.Details = string(sval)
			}
			// t.Bandwidth (int64) (int64)
		case "Bandwidth":
			{
				maj, extra, err := cr.ReadHeader()
				var extraI int64
				if err != nil {
					return err
				}
				switch maj {
				case cbg.MajUnsignedInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 positive overflow")
					}
				case cbg.MajNegativeInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 negative overflow")
					}
					extraI = -1 - extraI
				default:
					return fmt.Errorf("wrong type for int64 field: %d", maj)
				}

				t.Bandwidth = int64(extraI)
			}
			// t.Requester (string) (string)
		case "Requester":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.Requester = string(sval)
			}
			// t.RetryCount (int64) (int64)
		case "RetryCount":
			{
				maj, extra, err := cr.ReadHeader()
				var extraI int64
				if err != nil {
					return err
				}
				switch maj {
				case cbg.MajUnsignedInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 positive overflow")
					}
				case cbg.MajNegativeInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 negative overflow")
					}
					extraI = -1 - extraI
				default:
					return fmt.Errorf("wrong type for int64 field: %d", maj)
				}

				t.RetryCount = int64(extraI)
			}
			// t.EdgeReplicas (int64) (int64)
		case "EdgeReplicas":
			{
				maj, extra, err := cr.ReadHeader()
				var extraI int64
				if err != nil {
					return err
				}
				switch maj {
				case cbg.MajUnsignedInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 positive overflow")
					}
				case cbg.MajNegativeInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 negative overflow")
					}
					extraI = -1 - extraI
				default:
					return fmt.Errorf("wrong type for int64 field: %d", maj)
				}

				t.EdgeReplicas = int64(extraI)
			}
			// t.EdgeWaitings (int64) (int64)
		case "EdgeWaitings":
			{
				maj, extra, err := cr.ReadHeader()
				var extraI int64
				if err != nil {
					return err
				}
				switch maj {
				case cbg.MajUnsignedInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 positive overflow")
					}
				case cbg.MajNegativeInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 negative overflow")
					}
					extraI = -1 - extraI
				default:
					return fmt.Errorf("wrong type for int64 field: %d", maj)
				}

				t.EdgeWaitings = int64(extraI)
			}
			// t.CandidateReplicas (int64) (int64)
		case "CandidateReplicas":
			{
				maj, extra, err := cr.ReadHeader()
				var extraI int64
				if err != nil {
					return err
				}
				switch maj {
				case cbg.MajUnsignedInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 positive overflow")
					}
				case cbg.MajNegativeInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 negative overflow")
					}
					extraI = -1 - extraI
				default:
					return fmt.Errorf("wrong type for int64 field: %d", maj)
				}

				t.CandidateReplicas = int64(extraI)
			}
			// t.CandidateWaitings (int64) (int64)
		case "CandidateWaitings":
			{
				maj, extra, err := cr.ReadHeader()
				var extraI int64
				if err != nil {
					return err
				}
				switch maj {
				case cbg.MajUnsignedInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 positive overflow")
					}
				case cbg.MajNegativeInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 negative overflow")
					}
					extraI = -1 - extraI
				default:
					return fmt.Errorf("wrong type for int64 field: %d", maj)
				}

				t.CandidateWaitings = int64(extraI)
			}
			// t.ReplenishReplicas (int64) (int64)
		case "ReplenishReplicas":
			{
				maj, extra, err := cr.ReadHeader()
				var extraI int64
				if err != nil {
					return err
				}
				switch maj {
				case cbg.MajUnsignedInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 positive overflow")
					}
				case cbg.MajNegativeInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 negative overflow")
					}
					extraI = -1 - extraI
				default:
					return fmt.Errorf("wrong type for int64 field: %d", maj)
				}

				t.ReplenishReplicas = int64(extraI)
			}
			// t.EdgeReplicaSucceeds ([]string) (slice)
		case "EdgeReplicaSucceeds":

			maj, extra, err = cr.ReadHeader()
			if err != nil {
				return err
			}

			if extra > cbg.MaxLength {
				return fmt.Errorf("t.EdgeReplicaSucceeds: array too large (%d)", extra)
			}

			if maj != cbg.MajArray {
				return fmt.Errorf("expected cbor array")
			}

			if extra > 0 {
				t.EdgeReplicaSucceeds = make([]string, extra)
			}

			for i := 0; i < int(extra); i++ {

				{
					sval, err := cbg.ReadString(cr)
					if err != nil {
						return err
					}

					t.EdgeReplicaSucceeds[i] = string(sval)
				}
			}

			// t.CandidateReplicaSucceeds ([]string) (slice)
		case "CandidateReplicaSucceeds":

			maj, extra, err = cr.ReadHeader()
			if err != nil {
				return err
			}

			if extra > cbg.MaxLength {
				return fmt.Errorf("t.CandidateReplicaSucceeds: array too large (%d)", extra)
			}

			if maj != cbg.MajArray {
				return fmt.Errorf("expected cbor array")
			}

			if extra > 0 {
				t.CandidateReplicaSucceeds = make([]string, extra)
			}

			for i := 0; i < int(extra); i++ {

				{
					sval, err := cbg.ReadString(cr)
					if err != nil {
						return err
					}

					t.CandidateReplicaSucceeds[i] = string(sval)
				}
			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}
func (t *NodePulledResult) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write([]byte{164}); err != nil {
		return err
	}

	// t.Size (int64) (int64)
	if len("Size") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Size\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("Size"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Size")); err != nil {
		return err
	}

	if t.Size >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Size)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.Size-1)); err != nil {
			return err
		}
	}

	// t.NodeID (string) (string)
	if len("NodeID") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"NodeID\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("NodeID"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("NodeID")); err != nil {
		return err
	}

	if len(t.NodeID) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.NodeID was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.NodeID))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.NodeID)); err != nil {
		return err
	}

	// t.Status (int64) (int64)
	if len("Status") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Status\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("Status"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Status")); err != nil {
		return err
	}

	if t.Status >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Status)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.Status-1)); err != nil {
			return err
		}
	}

	// t.BlocksCount (int64) (int64)
	if len("BlocksCount") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"BlocksCount\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("BlocksCount"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("BlocksCount")); err != nil {
		return err
	}

	if t.BlocksCount >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.BlocksCount)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.BlocksCount-1)); err != nil {
			return err
		}
	}
	return nil
}

func (t *NodePulledResult) UnmarshalCBOR(r io.Reader) (err error) {
	*t = NodePulledResult{}

	cr := cbg.NewCborReader(r)

	maj, extra, err := cr.ReadHeader()
	if err != nil {
		return err
	}
	defer func() {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	if maj != cbg.MajMap {
		return fmt.Errorf("cbor input should be of type map")
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("NodePulledResult: map struct too large (%d)", extra)
	}

	var name string
	n := extra

	for i := uint64(0); i < n; i++ {

		{
			sval, err := cbg.ReadString(cr)
			if err != nil {
				return err
			}

			name = string(sval)
		}

		switch name {
		// t.Size (int64) (int64)
		case "Size":
			{
				maj, extra, err := cr.ReadHeader()
				var extraI int64
				if err != nil {
					return err
				}
				switch maj {
				case cbg.MajUnsignedInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 positive overflow")
					}
				case cbg.MajNegativeInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 negative overflow")
					}
					extraI = -1 - extraI
				default:
					return fmt.Errorf("wrong type for int64 field: %d", maj)
				}

				t.Size = int64(extraI)
			}
			// t.NodeID (string) (string)
		case "NodeID":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.NodeID = string(sval)
			}
			// t.Status (int64) (int64)
		case "Status":
			{
				maj, extra, err := cr.ReadHeader()
				var extraI int64
				if err != nil {
					return err
				}
				switch maj {
				case cbg.MajUnsignedInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 positive overflow")
					}
				case cbg.MajNegativeInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 negative overflow")
					}
					extraI = -1 - extraI
				default:
					return fmt.Errorf("wrong type for int64 field: %d", maj)
				}

				t.Status = int64(extraI)
			}
			// t.BlocksCount (int64) (int64)
		case "BlocksCount":
			{
				maj, extra, err := cr.ReadHeader()
				var extraI int64
				if err != nil {
					return err
				}
				switch maj {
				case cbg.MajUnsignedInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 positive overflow")
					}
				case cbg.MajNegativeInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 negative overflow")
					}
					extraI = -1 - extraI
				default:
					return fmt.Errorf("wrong type for int64 field: %d", maj)
				}

				t.BlocksCount = int64(extraI)
			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}
