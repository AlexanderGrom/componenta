package sqlx

// Компоненты запроса
type components struct {
	Aggregate []aggregateComponent
	Select    []interface{}
	Insert    []interface{}
	Update    []interface{}
	Delete    []interface{}
	From      []fromComponent
	Join      []joinComponent
	Into      []interface{}
	Columns   []interface{}
	Values    []valueComponent
	OrIgnore  []interface{}
	ReturnId  []interface{}
	Set       []setComponent
	Where     []whereComponent
	Group     []interface{}
	Having    []havingComponent
	Order     []orderComponent
	Limit     []interface{}
	Offset    []interface{}
}

func newComponents() *components {
	return &components{}
}

type aggregateComponent struct {
	function string
	column   interface{}
	alias    string
}

type fromComponent struct {
	kind    string
	table   interface{}
	builder *Builder
}

type joinComponent Joiner
type valueComponent []interface{}
type setComponent Data

type whereComponent struct {
	kind     string
	column   interface{}
	operator interface{}
	value    interface{}
	min      interface{}
	max      interface{}
	boolean  interface{}
	list     List
	builder  *Builder
}

type havingComponent struct {
	kind     string
	column   interface{}
	operator interface{}
	value    interface{}
	boolean  interface{}
	builder  *Builder
}

type orderComponent struct {
	column    string
	direction string
}
