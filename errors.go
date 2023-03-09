package bubbleviews

type EscPressedError struct{}

func (EscPressedError) Error() string { return "esc was pressed" }
