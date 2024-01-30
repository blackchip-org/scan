package sgo

import "github.com/blackchip-org/scan"

const (
	StringType = "string"
)

var (
	CharEnc = scan.NewCharEncRule(
		scan.AlertEnc,
		scan.BackspaceEnc,
		scan.FormFeedEnc,
		scan.LineFeedEnc,
		scan.CarriageReturnEnc,
		scan.HorizontalTabEnc,
		scan.VerticalTabEnc,
	)
	String = scan.NewStrRule('"', '"').
		WithType(StringType).
		WithEscape('\\').
		WithEscapeRules(
			CharEnc,
			scan.Hex2Enc,
			scan.Hex4Enc,
			scan.Hex8Enc,
			scan.OctEnc,
		)
)
