package style

// DefaultMaxOpts specifies the default maximum number of visible options/elements.
const DefaultMaxOpts = 5

// DefaultDistribution represents the standard horizontal layout distribution.
var DefaultDistribution = HorizontalDistribution(JustifyAround, DefaultMaxOpts)

// Distribution defines layout rules for organizing child elements along an axis.
type Distribution struct {
	// Direction specifies whether elements flow horizontally or vertically.
	Direction Direction
	// Justify controls how remaining space is distributed along the main axis.
	Justify Justify
	// Limit sets the maximum number of items allowed in this distribution layout.
	Limit uint16
}

// VerticalDistribution constructs a vertical layout distribution.
func VerticalDistribution() Distribution {
	return Distribution{
		Direction: Vertical,
	}
}

// HorizontalDistribution constructs a horizontal layout distribution with space justification and element limit.
func HorizontalDistribution(justify Justify, limit uint16) Distribution {
	return Distribution{
		Direction: Horizontal,
		Justify:   justify,
		Limit:     limit,
	}
}
