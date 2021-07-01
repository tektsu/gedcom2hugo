# gedcom2hugo

Generate Hugo content from a GEDCOM 5.5 file.

## Usage

```bash
$ go run . -gedcom /path/to/family.ged -project /path/to/target/
```

```bash
$ go build
$ ./gedcom2hugo -gedcom /path/to/family.ged -project /path/to/target/
```

## Contribution
Contribution is highly welcome! There is just one rule: use `go fmt` before committing. I don't want to discuss code
style, it's boring. There is a standard, follow it.

## Credits
This project is based on [the implementation](https://github.com/tektsu/gedcom2hugo) of [@tektsu](https://github.com/tektsu/)
done for his [idrisproject](https://www.idrisproject.com/) and uses the [gedcom parser](https://github.com/iand/gedcom)
from [@iand](https://github.com/iand/).
