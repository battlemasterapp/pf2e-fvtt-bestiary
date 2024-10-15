# PF2E Bestiary Data

This is a simple attempt to parse the foundryvtt PF2E compendium data into the format that I use for my applications. This is a work in progress and will be updated as I have time to work on it.

Instead of using individual files for each entry, I have decided to use a single file for each book. This is to reduce the number of files that need to be loaded and to make it easier to manage the data.

## Parser

Text descriptions can contain foundry commands. A parser cleans up the text and replaces the foundry commands with the appropriate content.

Commands can be identified by the following format:

```
@Command[param1|param2|...]{name}
```

Commands start with a `@` symbol and are followed by the command name. The command name is followed by a list of parameters separated by `|` characters. The command name and parameters are enclosed in square brackets. The `name` and the curly brackets are optional and used to replace the command with the `name` if present.

Available commands:

- `@UUID[id]{name}` - the parser replaces with the `name` if present. If not, it will use the last part of the `id`.
- `@Localize[id]` - the parser replaces with the localized string. The id is a key in the localization `static/lang/en.json` file. Example: id `foo.bar.baz` will be replaced with the value of `foo.bar.baz` in the localization file.
- `@Check[ability|dc:n|name|traits|basic]` - the parser replaces with a check description. The `ability` is the ability used for the check. The `dc` is the difficulty class. The `name` is the name of the check. The `traits` are the traits of the check. The `basic` is a flag that tell if the save is basic or not. `name`, `traits`, and `basic` are optional.
- `@Template[type|distance:n]{name}` - the parser replaces with the `name` if present. If not, it will use `n-foot type` where `n` is the distance.
- `@Damage[ndm[type]]` - the parser replaces with the damage description. `n` is the number of dice. `m` is the number of sides on the dice. `type` is the damage type. `type` is optional and if the type is `untyped` it's not displayed.