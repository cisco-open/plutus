# Parser

This pakcage is used to parse HCL policies.

#### Possible improvement

As of now, the parser treats the policies as HCL2 which causes the followin issue.

> Invalid block definition; A block definition must have block content delimited by \"{\" and \"}\", starting on the same line as the block header

The parser can be updated to make hte HCL version configurable.