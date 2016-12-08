package expvarmon

import (
	"strings"
	"testing"
	"time"

	"github.com/antonholmquist/jason"
)

func TestVarName(t *testing.T) {
	v := VarName("memstats.Alloc")

	slice := v.ToSlice()
	if len(slice) != 2 || slice[0] != "memstats" || slice[1] != "Alloc" {
		t.Fatalf("ToSlice failed: %v", slice)
	}

	short := v.Short()
	if short != "Alloc" {
		t.Fatalf("Expecting Short() to be 'Alloc', but got: %s", short)
	}

	kind := v.Kind()
	if kind != KindDefault {
		t.Fatalf("Expecting kind to be %v, but got: %v", KindDefault, kind)
	}

	v = VarName("mem:memstats.Alloc")

	slice = v.ToSlice()
	if len(slice) != 2 || slice[0] != "memstats" || slice[1] != "Alloc" {
		t.Fatalf("ToSlice failed: %v", slice)
	}

	short = v.Short()
	if short != "Alloc" {
		t.Fatalf("Expecting Short() to be 'Alloc', but got: %s", short)
	}

	kind = v.Kind()
	if kind != KindMemory {
		t.Fatalf("Expecting kind to be %v, but got: %v", KindMemory, kind)
	}

	v = VarName("duration:ResponseTimes.API.Users")
	kind = v.Kind()
	if kind != KindDuration {
		t.Fatalf("Expecting kind to be %v, but got: %v", KindDuration, kind)
	}

	// single \. escapes the dot
	v = VarName(`bleve.indexes.bench\.bleve.index.lookup_queue_len`)

	slice = v.ToSlice()
	if len(slice) != 5 || slice[0] != "bleve" || slice[1] != "indexes" || slice[2] != "bench.bleve" ||
		slice[3] != "index" || slice[4] != "lookup_queue_len" {
		t.Fatalf("ToSlice failed: %v", slice)
	}

	// double \\. escapes backslash, not dot
	v = VarName(`bleve.indexes.bench\\.bleve.index.lookup_queue_len`)

	slice = v.ToSlice()
	if len(slice) != 6 || slice[0] != "bleve" || slice[1] != "indexes" || slice[2] != "bench\\" ||
		slice[3] != "bleve" || slice[4] != "index" || slice[5] != "lookup_queue_len" {
		t.Fatalf("ToSlice failed: %v", slice)
	}

	// triple \\\. escapes backslash then dot
	v = VarName(`bleve.indexes.bench\\\.bleve.index.lookup_queue_len`)

	slice = v.ToSlice()
	if len(slice) != 5 || slice[0] != "bleve" || slice[1] != "indexes" || slice[2] != "bench\\.bleve" ||
		slice[3] != "index" || slice[4] != "lookup_queue_len" {
		t.Fatalf("ToSlice failed: %v", slice)
	}

	// quadruple \\\\. escapes two backslashes, not dot
	v = VarName(`bleve.indexes.bench\\\\.bleve.index.lookup_queue_len`)

	slice = v.ToSlice()
	if len(slice) != 6 || slice[0] != "bleve" || slice[1] != "indexes" || slice[2] != "bench\\\\" ||
		slice[3] != "bleve" || slice[4] != "index" || slice[5] != "lookup_queue_len" {
		t.Fatalf("ToSlice failed: %v", slice)
	}

	// unsupported \x passes through unaltered
	v = VarName(`bleve.indexes.bench\xbleve.index.lookup_queue_len`)

	slice = v.ToSlice()
	if len(slice) != 5 || slice[0] != "bleve" || slice[1] != "indexes" || slice[2] != "bench\\xbleve" ||
		slice[3] != "index" || slice[4] != "lookup_queue_len" {
		t.Fatalf("ToSlice failed: %v", slice)
	}
}

func TestVarNew(t *testing.T) {
	v := NewVar(VarName("mem:field.Subfield"))
	if v.Kind() != KindMemory {
		t.Fatalf("Expect Memory, got %v", v.Kind())
	}
	v = NewVar(VarName("str:field.Subfield"))
	if v.Kind() != KindString {
		t.Fatalf("Expect String, got %v", v.Kind())
	}
	v = NewVar(VarName("duration:field.Subfield"))
	if v.Kind() != KindDuration {
		t.Fatalf("Expect Duration, got %v", v.Kind())
	}
}

func str2val(t *testing.T, s string) *jason.Value {
	val, err := jason.NewValueFromReader(strings.NewReader(s))
	if err != nil {
		t.Fatal(err)
	}
	return val
}

func TestVarNumber(t *testing.T) {
	v := &Number{}
	testNumber := func(t *testing.T, v *Number, json string, intval int, str string) {
		v.Set(str2val(t, json))
		if want := intval; v.Value() != want {
			t.Fatalf("Expect value to be %d, got %d", want, v.Value())
		}
		if want := str; v.String() != want {
			t.Fatalf("Expect value to be %s, got %s", want, v.String())
		}
	}

	testNumber(t, v, "142", 142, "142")
	testNumber(t, v, "13.24", 13, "13.24")
	testNumber(t, v, "true", 0, "0")
	testNumber(t, v, "\"success\"", 0, "0")
}

func TestVarMemory(t *testing.T) {
	v := &Memory{}
	testMemory := func(t *testing.T, v *Memory, json string, intval int, str string) {
		v.Set(str2val(t, json))
		if want := intval; v.Value() != want {
			t.Fatalf("Expect value to be %d, got %d", want, v.Value())
		}
		if want := str; v.String() != want {
			t.Fatalf("Expect value to be %s, got %s", want, v.String())
		}
	}

	testMemory(t, v, "12", 12, "12B")
	testMemory(t, v, "1024", 1024, "1.0KB")
	testMemory(t, v, "1048576", 1048576, "1.0MB")
	testMemory(t, v, "1073741824", 1073741824, "1.0GB")
	testMemory(t, v, "6815744", 6815744, "6.5MB")
	testMemory(t, v, "128849018880", 128849018880, "120GB")
}

func TestVarDuration(t *testing.T) {
	v := &Duration{}
	testDuration := func(t *testing.T, v *Duration, json string, intval int, str string) {
		v.Set(str2val(t, json))
		if want := intval; v.Value() != want {
			t.Fatalf("Expect value to be %d, got %d", want, v.Value())
		}
		if want := str; v.String() != want {
			t.Fatalf("Expect value to be %s, got %s", want, v.String())
		}
	}

	testDuration(t, v, "12", 12, "12ns")
	testDuration(t, v, "1000", 1e3, "1µs")
	testDuration(t, v, "2000", 2*1e3, "2µs")
	testDuration(t, v, "1000000", 1e6, "1ms")
	testDuration(t, v, "2000000", 2*1e6, "2ms")
	testDuration(t, v, "155000000", 155*1e6, "155ms")
	testDuration(t, v, "1000000000", 1e9, "1s")
	testDuration(t, v, "13000000000", 13*1e9, "13s")
	testDuration(t, v, "60000000000", 60*1e9, "1m0s")
	testDuration(t, v, "90000000000", 90*1e9, "1m30s")
	testDuration(t, v, "172800000000000", 48*3600*1e9, "48h0m0s")
	testDuration(t, v, "63072000000000000", 2*365*24*3600*1e9, "17520h0m0s")
}

func TestVarString(t *testing.T) {
	v := &String{}
	v.Set(str2val(t, "\"success\""))
	if want := "success"; v.String() != want {
		t.Fatalf("Expect value to be %s, got %s", want, v.String())
	}
	v.Set(str2val(t, "123"))
	if want := "N/A"; v.String() != want {
		t.Fatalf("Expect value to be %s, got %s", want, v.String())
	}
}

func TestVarGCPauses(t *testing.T) {
	v := &GCPauses{}
	v.Set(str2val(t, pauseNSTest))
	hist := v.Histogram(20)
	if values, _ := hist.BarchartData(); len(values) != 20 {
		t.Fatalf("Expect len of values to be 20, got %v", len(values))
	}
	// TODO: check if it's true mean :)
	if want, mean := 112746.86928104576, hist.Mean(); mean != want {
		t.Fatalf("Expect mean to be be %v, got %v", want, mean)
	}
}

func TestVarGCIntervals(t *testing.T) {
	v := &GCIntervals{}
	v.Set(str2val(t, pauseEndTest))
	hist := v.Histogram(20)
	if values, _ := hist.BarchartData(); len(values) != 20 {
		t.Fatalf("Expect len of values to be 20, got %v", len(values))
	}
	// TODO: check if it's true mean :)
	if want, mean := time.Duration(125205658), time.Duration(hist.Mean()); mean != want {
		t.Fatalf("Expect mean to be be %v, got %v", want, mean)
	}

	v.Set(str2val(t, pauseEnd2Test))
	hist = v.Histogram(20)
	if want, mean := time.Duration(0), time.Duration(hist.Mean()); mean != want {
		t.Fatalf("Expect mean to be be %v, got %v", want, mean)
	}
	if want, max := time.Duration(0), time.Duration(hist.Max()); max != want {
		t.Fatalf("Expect max to be be %v, got %v", want, max)
	}

	v.Set(str2val(t, pauseEnd3Test))
	hist = v.Histogram(20)
	if want, min := time.Duration(100), time.Duration(hist.Min()); min != want {
		t.Fatalf("Expect min to be be %v, got %v", want, min)
	}
	if want, mean := time.Duration(100), time.Duration(hist.Mean()); mean != want {
		t.Fatalf("Expect mean to be be %v, got %v", want, mean)
	}
	if want, max := time.Duration(100), time.Duration(hist.Max()); max != want {
		t.Fatalf("Expect max to be be %v, got %v", want, max)
	}

	v.Set(str2val(t, pauseEnd4Test))
	hist = v.Histogram(20)
	if want, min := time.Duration(50), time.Duration(hist.Min()); min != want {
		t.Fatalf("Expect min to be be %v, got %v", want, min)
	}
	if want, mean := time.Duration(100), time.Duration(hist.Mean()); mean != want {
		t.Fatalf("Expect mean to be be %v, got %v", want, mean)
	}
	if want, p95 := time.Duration(150), time.Duration(hist.Quantile(0.99)); p95 != want {
		t.Fatalf("Expect mean to be be %v, got %v", want, p95)
	}
	if want, max := time.Duration(150), time.Duration(hist.Max()); max != want {
		t.Fatalf("Expect max to be be %v, got %v", want, max)
	}
}

const (
	pauseNSTest   = "[65916, 92412, 67016, 59076, 55161, 53428, 128675, 90476, 78093, 60473, 64353, 58214, 83926, 64390, 103391, 71275, 76651, 56475, 367180, 184505, 307648, 175680, 129120, 102616, 127322, 224862, 83092, 148607, 122833, 139011, 494885, 97452, 95129, 115403, 119657, 122214, 111744, 115824, 95834, 81927, 91120, 131541, 75511, 135424, 125637, 85784, 107094, 101551, 110081, 80628, 123030, 130343, 128940, 114670, 111470, 75146, 101250, 117553, 112062, 106360, 101543, 108607, 245857, 106147, 108091, 84570, 78700, 117863, 74284, 102977, 83952, 108068, 89709, 115250, 108062, 135150, 84460, 389962, 109881, 79255, 88669, 106366, 90551, 115548, 93409, 124459, 93660, 132709, 70662, 119209, 86984, 118776, 114768, 107875, 70117, 95590, 90558, 86439, 85069, 83155, 89212, 115581, 61221, 78387, 67468, 82099, 107160, 83947, 109817, 113753, 121822, 87682, 104144, 88659, 82247, 91591, 138847, 498527, 121882, 114585, 135840, 111263, 101143, 106915, 100841, 110974, 71145, 97220, 118328, 103716, 115043, 74672, 86126, 106929, 115845, 97969, 118960, 103949, 96019, 80543, 106717, 115346, 114901, 88455, 76337, 107155, 141398, 92871, 120444, 90579, 110057, 94518, 115869, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]"
	pauseEndTest  = "[1479043813714676394, 1479043813826878234, 1479043813945138433, 1479043814068508212, 1479043814189654357, 1479043814308676597, 1479043814431564320, 1479043814562593715, 1479043814711247081, 1479043814834144113, 1479043814985132450, 1479043815100092005, 1479043815212510124, 1479043815339516264, 1479043815466954803, 1479043815586638242, 1479043815701466921, 1479043815822886110, 1479043815940787320, 1479043816057204555, 1479043816172730959, 1479043816290981742, 1479043816420391704, 1479043816536639680, 1479043816667774007, 1479043816781354036, 1479043816914947573, 1479043817032592375, 1479043817151428431, 1479043817315778541, 1479043817446275080, 1479043817567438718, 1479043817690605256, 1479043817818987622, 1479043817941235852, 1479043818066955733, 1479043818187269104, 1479043818314026636, 1479043818427935634, 1479043818551454011, 1479043818699543272, 1479043818809199581, 1479043818923638543, 1479043819037011907, 1479043819163217842, 1479043819290042678, 1479043819418901787, 1479043819533917812, 1479043819654132939, 1479043819783596738, 1479043819891910663, 1479043820026240593, 1479043820146020373, 1479043820285220624, 1479043820412476764, 1479043820536800828, 1479043820668164426, 1479043820792675685, 1479043820912623250, 1479043821037083034, 1479043821173056176, 1479043821291083992, 1479043821410303816, 1479043821534180082, 1479043821720749924, 1479043821841731195, 1479043821955876417, 1479043822078035785, 1479043822203369961, 1479043822326472365, 1479043822459803872, 1479043822566177411, 1479043822687468874, 1479043822806338268, 1479043822920227164, 1479043823038173722, 1479043823151756215, 1479043823272830087, 1479043823405407596, 1479043823532679112, 1479043823669829071, 1479043823778213707, 1479043823895712316, 1479043824010952001, 1479043824124393414, 1479043824245079718, 1479043824374113198, 1479043824495812676, 1479043824610975198, 1479043824721213448, 1479043824840472825, 1479043824983362919, 1479043825181166123, 1479043793253721601, 1479043793385360442, 1479043793502007593, 1479043793639226256, 1479043793772090554, 1479043793899206006, 1479043794023356018, 1479043794141100539, 1479043794277195410, 1479043794398575041, 1479043794520214726, 1479043794626326933, 1479043794742733618, 1479043794860765941, 1479043794989443603, 1479043795115450686, 1479043795239084332, 1479043795368101533, 1479043795479290335, 1479043795617320435, 1479043795738300857, 1479043795853773091, 1479043795986557268, 1479043796105841750, 1479043796217019850, 1479043796338541232, 1479043796464565411, 1479043796578884858, 1479043796707472399, 1479043796830277286, 1479043796953986052, 1479043797068021156, 1479043797186576086, 1479043797310551203, 1479043797444077233, 1479043797569716319, 1479043797692845425, 1479043797817668329, 1479043797933388910, 1479043798045590714, 1479043798156713161, 1479043798306287956, 1479043798412917001, 1479043798534373963, 1479043798672883517, 1479043798796160178, 1479043798935747512, 1479043799052631151, 1479043799170142721, 1479043799314178880, 1479043799429953705, 1479043799574795360, 1479043799680393451, 1479043799809439283, 1479043799940976218, 1479043800092305564, 1479043800210178253, 1479043800335567367, 1479043800461761677, 1479043800578874973, 1479043800693208302, 1479043800814372466, 1479043800930977672, 1479043801072086934, 1479043801205701809, 1479043801347696907, 1479043801475107270, 1479043801593208254, 1479043801718445855, 1479043801859340806, 1479043802005962046, 1479043802129058466, 1479043802242484743, 1479043802373364585, 1479043802511716704, 1479043802646728106, 1479043802784023364, 1479043802915246996, 1479043803044801734, 1479043803196646614, 1479043803324345525, 1479043803443686755, 1479043803562811189, 1479043803694057521, 1479043803808555206, 1479043803928702376, 1479043804046162087, 1479043804170210698, 1479043804298680408, 1479043804420335262, 1479043804541609195, 1479043804663594548, 1479043804775065999, 1479043804910592345, 1479043805045321238, 1479043805159266647, 1479043805289679498, 1479043805405792163, 1479043805525491606, 1479043805640852468, 1479043805762590003, 1479043805881520640, 1479043806004392626, 1479043806125713978, 1479043806248657803, 1479043806376939748, 1479043806500327139, 1479043806619621564, 1479043806741777582, 1479043806875046280, 1479043806993319258, 1479043807133805788, 1479043807254364828, 1479043807374886160, 1479043807491091944, 1479043807603905843, 1479043807745406015, 1479043807875727264, 1479043808009547001, 1479043808136933791, 1479043808257632715, 1479043808390160950, 1479043808500101804, 1479043808625577483, 1479043808740897854, 1479043808856659338, 1479043808985295029, 1479043809115161228, 1479043809219124395, 1479043809361693329, 1479043809480624039, 1479043809603959486, 1479043809714056133, 1479043809828315044, 1479043809952450826, 1479043810124934166, 1479043810251188760, 1479043810368719081, 1479043810486584295, 1479043810605792657, 1479043810723171767, 1479043810834296984, 1479043810964120950, 1479043811125530295, 1479043811251944952, 1479043811372830005, 1479043811548218962, 1479043811681514638, 1479043811807009040, 1479043811941029199, 1479043812069834652, 1479043812185333792, 1479043812305107561, 1479043812429889877, 1479043812556417122, 1479043812681336827, 1479043812804958914, 1479043812936661566, 1479043813054411177, 1479043813174215463, 1479043813348933717, 1479043813459811570, 1479043813582526620]"
	pauseEnd2Test = "[1479053221958038493, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]"
	pauseEnd3Test = "[1479053221958038400, 1479053221958038500, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]"
	pauseEnd4Test = "[1479053221958038400, 1479053221958038450, 1479053221958038550, 1479053221958038700, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]"
)
