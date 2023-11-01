# How we document our command line syntax

## Literal text

Use plain text for parts of the command that cannot be changed.

_example:_
`goctl help`
The argument help is required in this command.

## Placeholder values

Use angled brackets to represent a value the user must replace. No other expressions can be contained within the angled brackets.

_example:_
`goctl pr view <issue-number>`
Replace `<issue-number>` with an issue number.

## Optional arguments

Place optional arguments in square brackets. Mutually exclusive arguments can be included inside square brackets if they are separated with vertical bars.

_example:_
`goctl pr checkout [--web]`
The argument `--web` is optional.

`goctl pr view [<number> | <url>]`
The `<number>` and `<url>` arguments are optional.

## Required mutually exclusive arguments

Place required mutually exclusive arguments inside braces, separate arguments with vertical bars.

_example:_
`goctl pr {view | create}`

## Repeatable arguments

Ellipsis represent arguments that can appear multiple times.

_example:_
`goctl pr close <pr-number>...`

## Variable naming

For multi-word variables use dash-case (all lower case with words separated by dashes)

_example:_
`goctl pr checkout <issue-number>`

## Additional examples

_optional argument with placeholder:_
`command sub-command [<arg>]`

_required argument with mutually exclusive options:_
`command sub-command {<path> | <string> | literal}`

_optional argument with mutually exclusive options:_
`command sub-command [<path> | <string>]`
