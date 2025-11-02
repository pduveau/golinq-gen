# OpenXML Linq to go structures

This tool parses the linq html files from https://docs.dndocs.com/n/DocumentFormat.OpenXml.Linq/3.1.0/api/ and generates an FS structure
containing the go packages and classes to manage xml document based on this structure.

A basic API to ```set```, ```add``` and ```get``` ```attributes``` and ```children``` is generated within the packages.

# How circular dependencies are managed

By default, each namespace is implemented as a go package. But some of them have circular dependencies.

To handle this constraints, two strategies are used:
- some namespaces are grouped in the most generic one:
  - when related to versions of the same concept (i.e.: ```a``` and ```a14```),
  - when the links are multiples
- using dynamic links instead of static ones:
  - two classes ```ext``` and ```graphicData``` of namespace ```a``` are extracted in a dedicated package that imports almost all of others
  - they are replaced by a simple class that export required method
  - the dynamic link is then established while calling ```linq.initLinq()``` prior to use the go structure.

# How to use

The package ```github.com/pduveau/go-office``` contains API designed to interact with this module in order to manage office documents.

The dynamic link requires to use two functions to switch the object from the interfaces to the real types:
- ```linq.Ext_a_to_Ext_aacircular``` convert a ```a.Ext_a``` to the ```aacircular.Ext_a``` which exposes all the method to manage the xml object &lt;a:ext&gt;.
- ```linq.GraphicData_a_to_GraphicData_aacircular``` convert a ```a.GraphicData_a``` to the ```aacircular.GraphicData_a``` which exposes all the method to manage the xml object &lt;a:graphicData&gt;.

The command line is: 
```bash
golinq-gen [-folder <folder>] [-github <github>] [-help]
```
Option:
- ```<folder>``` is the destination folder. Default ```../go-office/```
- ```<github>``` is the github root of the target module. Default ```github.com/pduveau/go-office```

If you would like to contribute to ```github.com/pduveau/go-office``` then clone it and generate the linq structure in it keeping the default github

Avoid to generate in the current folder, as this tool create a folder where it stores the source html files in order to cache them.