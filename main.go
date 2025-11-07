package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

//go:embed templates/common/types.go
var types []byte

//go:embed templates/common/empty.go
var empty []byte

//go:embed templates/common/val.go
var val []byte

//go:embed templates/common/text.go
var text []byte

var github = "github.com/pduveau/go-office"

func main() {
	var f *os.File
	var err error
	folder := filepath.Join("..", "go-office")

	err = loadTemplates()

	flag.StringVar(&github, "github", github, "Specify the github repository to fill.")
	flag.StringVar(&folder, "folder", folder, "Specify the destination package folder.")
	help := flag.Bool("help", false, "This help.")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	folder = filepath.Join(folder, "linq")

	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	urls := getToc()

	for _, u := range urls {
		f := getDataReader(u)
		parseNamespace(f, u)
	}

	patchA()
	patchRId()
	patchMC()
	patchW()

	patchText()

	prepare()

	// delete(xmlElements["w"]["rPr"].Children, "w:shadow")

	common := filepath.Join(folder, "common")
	err = os.MkdirAll(common, 0775)

	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	aacircular := filepath.Join(folder, "aacircular")
	err = os.MkdirAll(aacircular, 0775)
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	f, err = os.Create(filepath.Join(common, "types.go"))
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	f.Write(types)
	f.Close()

	f, err = os.Create(filepath.Join(common, "empty.go"))
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	f.Write(empty)
	f.Close()

	f, err = os.Create(filepath.Join(common, "val.go"))
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	f.Write(val)
	f.Close()

	f, err = os.Create(filepath.Join(common, "text.go"))
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	f.Write(text)
	f.Close()

	parseInitLinq(folder)

	for pkg, classes := range xmlElements {
		if pkg == "" || len(classes) == 0 {
			continue
		}

		for n, cls := range classes {
			dir := filepath.Join(folder, cls.Element.Gopackage)
			outofns := ""
			if pkg == "a" && (n == "ext" || n == "graphicData") {
				cls.createDerivedFile(dir)
				outofns = "aacircular"
				dir = aacircular
			}

			err = cls.createClassFile(dir, outofns)
			if err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(1)
			}
		}
	}
}

func patchText() {
	type list struct {
		text    []string
		raw     []string
		exclude []string
		hasText []string
	}

	apply := map[string]list{
		"a": {
			text: []string{"tableStyleId", "t"},
		},
		"ap": {
			exclude: []string{"DigSig", "HeadingPairs", "HLinks", "Properties", "TitlesOfParts"},
		},
		"ask": {
			text: []string{"seed"},
		},
		"b": {
			exclude: []string{"Artist", "Author", "BookAuthor", "Compiler", "Composer", "Conductor", "Counsel",
				"Director", "Editor", "Interviewee", "Interviewer", "Inventor", "NameList", "Performer", "Person",
				"ProducerName", "Source", "Sources", "Translator", "Writer"},
		},
		"c": {
			text: []string{"evenFooter", "evenHeader", "firstFooter", "firstHeader", "formatCode", "f",
				"name", "oddFooter", "oddHeader", "separator", "v"},
		},
		"c15": {
			text: []string{"f", "sqref", "txfldGUID"},
		},
		"cdr": {
			text: []string{"x", "y"},
		},
		"cs": {
			text: []string{"lineWidthScale"},
		},
		"cx": {
			text: []string{"axisId", "binary", "binCount", "binSize", "copyright", "entityType", "evenFooter",
				"evenHeader", "f", "firstFooter", "firstHeader", "oddFooter", "oddHeader", "nf", "separator", "v"},
			hasText: []string{"idx", "pt"},
		},
		"emma": {
			text: []string{"literal"},
		},
		"inkml": {
			hasText: []string{"annotation", "matrix", "table", "trace"},
		},
		"m": {
			hasText: []string{"t"},
		},
		"msink": {
			hasText: []string{"property"},
		},
		"o": {
			text: []string{"FieldCodes", "LinkType", "LockedField"},
		},
		"oac": {
			text: []string{"imgData", "origImgData"},
		},
		"p": {
			text: []string{"attrName", "text"},
		},
		"vt": {
			exclude: []string{"array", "vector", "empty", "null", "variant"},
			hasText: []string{"vstream", "cf"},
		},
		"w": {
			raw:  []string{"instrText", "delInstrText"},
			text: []string{"delText", "fldData", "t"},
		},
		"wne": {
			text: []string{"eventDocBuildingBlockAfterInsert", "eventDocClose", "eventDocContentControlAfterInsert",
				"eventDocContentControlBeforeDelete", "eventDocContentControlContentUpdate", "eventDocContentControlOnEnter",
				"eventDocContentControlOnExit", "eventDocNew", "eventDocOpen", "eventDocStoreUpdate",
				"eventDocSync", "eventDocXmlAfterInsert", "eventDocXmlBeforeDelete"},
		},
		"wp": {
			text: []string{"align", "posOffset"},
		},
		"wp14": {
			text: []string{"pctHeight", "pctPosHOffset", "pctPosVOffset", "pctWidth"},
		},
		"x": {
			hasText: []string{"definedName", "f", "t"},
			text: []string{"author", "calculatedColumnFormula", "evenFooter", "evenHeader", "firstFooter", "firstHeader", "formula",
				"formula1", "formula2", "oddFooter", "oddHeader", "oldFormula", "stp", "totalsRowFormula", "v", "val"},
		},
		"x12ac": {
			text: []string{"list"},
		},
		"x14": {
			hasText: []string{"argumentDescription", "editValue"},
			text:    []string{"id", "tupleItem"},
		},
		"x15": {
			text: []string{"v"},
		},
		"xdr": {
			text: []string{"col", "colOff", "row", "rowOff"},
		},
		"xlrd": {
			hasText: []string{"fb"},
			text:    []string{"v"},
		},
		"xlrd2": {
			hasText: []string{"v", "rpv"},
		},
		"xltc": {
			text: []string{"text"},
		},
		"xne": {
			text: []string{"f", "sqref"},
		},
		"xvml": {
			exclude: []string{"ClientData"},
		},
		"xxpim": {
			text: []string{"implicitMeasureSupport"},
		},
		"xxpvi": {
			text: []string{"lastRefreshFeature", "lastUpdateFeature", "requiredFeature"},
		},
	}

	for n, t := range apply {
		for _, c := range t.text {
			if _, ok := xmlElements[n][c]; ok {
				xmlElements[n][c].classType = Text
			} else {
				fmt.Printf("Text %s:%s don't exist !\n", n, c)
			}
		}
		for _, c := range t.raw {
			if _, ok := xmlElements[n][c]; ok {
				xmlElements[n][c].classType = RawText
			} else {
				fmt.Printf("Raw %s:%s don't exist !\n", n, c)
			}
		}
		for _, c := range t.hasText {
			if _, ok := xmlElements[n][c]; ok {
				xmlElements[n][c].classType = HasText
			} else {
				fmt.Printf("hasText %s:%s don't exist !\n", n, c)
			}
		}
		if len(t.exclude) > 0 {
			for _, c := range t.exclude {
				if _, ok := xmlElements[n][c]; !ok {
					fmt.Printf("Exclude %s:%s don't exist !\n", n, c)
				}
			}
			for c := range xmlElements[n] {
				if !slices.Contains(t.exclude, c) && !slices.Contains(t.raw, c) && !slices.Contains(t.hasText, c) {
					xmlElements[n][c].classType = Text
				}
			}
		}
	}

}

func patchRId() {
	xmlAttributes["r"]["id"].Goname = "RId"
}

func patchA() {
	xmlElements["a"]["fillRect"].Xmlattribs = append(xmlElements["a"]["fillRect"].Xmlattribs,
		xmldata{Name: "b", Url: "DocumentFormat.OpenXml.Linq.NoNamespace.html"},
		xmldata{Name: "l", Url: "DocumentFormat.OpenXml.Linq.NoNamespace.html"},
		xmldata{Name: "t", Url: "DocumentFormat.OpenXml.Linq.NoNamespace.html"},
		xmldata{Name: "r", Url: "DocumentFormat.OpenXml.Linq.NoNamespace.html"},
	)
	xmlElements["a"]["ext"].Xmlattribs = append(xmlElements["a"]["ext"].Xmlattribs,
		xmldata{Name: "cx", Url: "DocumentFormat.OpenXml.Linq.A.html"},
		xmldata{Name: "cy", Url: "DocumentFormat.OpenXml.Linq.A.html"},
	)
	xmlAttributes["a"]["cx"] = &element{
		Goname: "Cx",
		ELocal: "cx",
		ETag:   "a:cx",
	}
	xmlAttributes["a"]["cy"] = &element{
		Goname: "Cy",
		ELocal: "cy",
		ETag:   "a:cy",
	}
}

func patchW() {
	b := &class{
		Aliases: []string{"bottom", "left", "right", "top", "insideH", "insideV", "tl2br", "tr2bl", "end"},
		Element: element{
			Goname:      "Border",
			GonameShort: "Border",
			ELocal:      "border",
			ETag:        "*",
		},
		Xmlattribs: []xmldata{
			{Name: "color", Url: "DocumentFormat.OpenXml.Linq.W.html"},
			{Name: "frame", Url: "DocumentFormat.OpenXml.Linq.W.html"},
			{Name: "shadow", Url: "DocumentFormat.OpenXml.Linq.W.html"},
			{Name: "size", Url: "DocumentFormat.OpenXml.Linq.W.html"},
			{Name: "themeColor", Url: "DocumentFormat.OpenXml.Linq.W.html"},
			{Name: "themeShade", Url: "DocumentFormat.OpenXml.Linq.W.html"},
			{Name: "themeTint", Url: "DocumentFormat.OpenXml.Linq.W.html"},
			{Name: "val", Url: "DocumentFormat.OpenXml.Linq.W.html"},
		},
		Children:   make(map[string]*class),
		Attributes: make(map[string]*element),
		classType:  Shared,
	}
	xmlElements["w"]["Border"] = b
	xmlAliases["w"] = map[string]*class{
		"bottom":  b,
		"left":    b,
		"right":   b,
		"top":     b,
		"insideH": b,
		"insideV": b,
		"tl2br":   b,
		"tr2bl":   b,
		"end":     b,
	}

	delete(xmlElements["w"], "top")
	delete(xmlElements["w"], "bottom")
	delete(xmlElements["w"], "end")
	delete(xmlElements["w"], "left")
	delete(xmlElements["w"], "right")
	delete(xmlElements["w"], "insideH")
	delete(xmlElements["w"], "insideV")
	delete(xmlElements["w"], "tl2br")
	delete(xmlElements["w"], "tr2bl")

	xmlElements["w"]["tblW"].Xmlattribs = append(xmlElements["w"]["tblW"].Xmlattribs,
		xmldata{Name: "w", Url: "DocumentFormat.OpenXml.Linq.W.html"},
		xmldata{Name: "type", Url: "DocumentFormat.OpenXml.Linq.W.html"},
	)

	xmlAttributes["w"]["w"] = &element{
		Goname: "W",
		ELocal: "w",
		ETag:   "w:w",
	}

	xmlAttributes["w"]["frame"] = &element{
		Goname: "Frame",
		ELocal: "frame",
		ETag:   "w:frame",
	}
	xmlAttributes["w"]["shadow"] = &element{
		Goname: "Shadow",
		ELocal: "shadow",
		ETag:   "w:shadow",
	}
	xmlAttributes["w"]["size"] = &element{
		Goname: "Size",
		ELocal: "size",
		ETag:   "w:size",
	}

	for _, p := range []string{"basedOn", "pStyle", "rStyle", "numId", "sz", "szCs", "lang"} {
		xmlElements["w"][p].classType = Valclass
	}

	xmlElements["w"]["document"].isProperties = true
	xmlElements["w"]["hyperlink"].isProperties = true
}

func patchMC() {
	xmlElements["mc"] = make(map[string]*class)
	xmlElements["mc"]["AlternateContent"] = &class{
		Aliases: []string{},
		Element: element{
			Goname:      "AlternateContent",
			GonameShort: "AlternateContent",
			ELocal:      "AlternateContent",
			ETag:        "mc:AlternateContent",
		},
		Xmlchildren: []xmldata{
			{Name: "Choice", Url: "DocumentFormat.OpenXml.Linq.MC.html"},
			{Name: "Fallback", Url: "DocumentFormat.OpenXml.Linq.MC.html"},
		},
		Xmlattribs: []xmldata{
			{Name: "Ignorable", Url: "DocumentFormat.OpenXml.Linq.MC.html"},
			{Name: "MustUnderstand", Url: "DocumentFormat.OpenXml.Linq.MC.html"},
			{Name: "ProcessContent", Url: "DocumentFormat.OpenXml.Linq.MC.html"},
			{Name: "PreserveElements", Url: "DocumentFormat.OpenXml.Linq.MC.html"},
			{Name: "PreserveAttributes", Url: "DocumentFormat.OpenXml.Linq.MC.html"},
		},
		Children:   make(map[string]*class),
		Attributes: make(map[string]*element),
	}
	xmlElements["mc"]["Choice"] = &class{
		Aliases: []string{},
		Element: element{
			Goname:      "Choice",
			GonameShort: "Choice",
			ELocal:      "Choice",
			ETag:        "mc:Choice",
		},
		Xmlchildren: []xmldata{
			{Name: "*", Url: "*"},
		},
		Xmlattribs: []xmldata{
			{Name: "Ignorable", Url: "DocumentFormat.OpenXml.Linq.MC.html"},
			{Name: "MustUnderstand", Url: "DocumentFormat.OpenXml.Linq.MC.html"},
			{Name: "ProcessContent", Url: "DocumentFormat.OpenXml.Linq.MC.html"},
			{Name: "PreserveElements", Url: "DocumentFormat.OpenXml.Linq.MC.html"},
			{Name: "PreserveAttributes", Url: "DocumentFormat.OpenXml.Linq.MC.html"},
			{Name: "Requires", Url: "DocumentFormat.OpenXml.Linq.MC.html"},
		},
		Children:   make(map[string]*class),
		Attributes: make(map[string]*element),
	}
	xmlElements["mc"]["Fallback"] = &class{
		Aliases: []string{},
		Element: element{
			Goname:      "Fallback",
			GonameShort: "Fallback",
			ELocal:      "Fallback",
			ETag:        "mc:Fallback",
		},
		Xmlchildren: []xmldata{
			{Name: "*", Url: "*"},
		},
		Xmlattribs: []xmldata{
			{Name: "Ignorable", Url: "DocumentFormat.OpenXml.Linq.MC.html"},
			{Name: "MustUnderstand", Url: "DocumentFormat.OpenXml.Linq.MC.html"},
			{Name: "ProcessContent", Url: "DocumentFormat.OpenXml.Linq.MC.html"},
			{Name: "PreserveElements", Url: "DocumentFormat.OpenXml.Linq.MC.html"},
			{Name: "PreserveAttributes", Url: "DocumentFormat.OpenXml.Linq.MC.html"},
		},
		Children:   make(map[string]*class),
		Attributes: make(map[string]*element),
	}

	xmlAttributes["mc"] = map[string]*element{}
	xmlAttributes["mc"]["Ignorable"] = &element{
		Goname: "Ignorable",
		ELocal: "Ignorable",
		ETag:   "Ignorable",
	}
	xmlAttributes["mc"]["MustUnderstand"] = &element{
		Goname: "MustUnderstand",
		ELocal: "MustUnderstand",
		ETag:   "MustUnderstand",
	}
	xmlAttributes["mc"]["ProcessContent"] = &element{
		Goname: "ProcessContent",
		ELocal: "ProcessContent",
		ETag:   "ProcessContent",
	}
	xmlAttributes["mc"]["PreserveElements"] = &element{
		Goname: "PreserveElements",
		ELocal: "PreserveElements",
		ETag:   "PreserveElements",
	}
	xmlAttributes["mc"]["PreserveAttributes"] = &element{
		Goname: "PreserveAttributes",
		ELocal: "PreserveAttributes",
		ETag:   "PreserveAttributes",
	}
	xmlAttributes["mc"]["Requires"] = &element{
		Goname: "Requires",
		ELocal: "Requires",
		ETag:   "Requires",
	}

	url2namespace["DocumentFormat.OpenXml.Linq.MC.html"] = "mc"
}
