## Mapping des CustomClasses C# vers Go

Ce projet contient une traduction progressive des classes C# du dossier `CustomClasses`
vers des packages Go internes. Le but est de conserver des noms proches tout en
respectant les conventions idiomatiques de Go.

- `CustomClasses/Save/DataFormat/*` → `internal/save/dataformat`
  - `SCS_Color` → `dataformat.Color`
  - `SCS_Float` → `dataformat.Float`
  - `SCS_string` → `dataformat.String`
  - `Vector_2f` → `dataformat.Vector2f`
  - `Vector_3f` → `dataformat.Vector3f`
  - `Vector_4f` / `Vector_3f_4f` → `dataformat.Vector4f`
  - `Vector_3i` → `dataformat.Vector3i`

- `CustomClasses/Save/Items/*` → `internal/save/items`
  - `Bank_Loan` → `items.BankLoan`
  - `SiiNBlockCore` → `items.SiiNBlockCore`
  - Les autres classes (`Bank`, `Garage`, `Player`, `Trailer`, etc.) ont chacune
    un équivalent Go portant le même nom en CamelCase (par ex. `Garage` → `items.Garage`).

- `CustomClasses/Save/ItemsExtra/*` → `internal/save/itemsextra`
  - `Cargo` → `itemsextra.Cargo`
  - `City` → `itemsextra.City`
  - `Company` → `itemsextra.Company`
  - `CompanyTruck` → `itemsextra.CompanyTruck`
  - `Country` → `itemsextra.Country`
  - `Garages` → `itemsextra.Garages`
  - `PlayerJob` → `itemsextra.PlayerJob`
  - `TrailerDefinition` → `itemsextra.TrailerDefinition`
  - `UserCompany*` → `itemsextra.UserCompany...`
  - `VisitedCity` → `itemsextra.VisitedCity`

- `CustomClasses/Save/Info/Dependencies.cs` → `internal/save/info`
  - `Dependencies` → `info.Dependencies`

- `CustomClasses/Save/SaveFileInfoData.cs` / `SaveFileProfileData.cs` → `internal/save`
  - `SaveFileInfoData` → `save.FileInfoData`
  - `SaveFileProfileData` → `save.FileProfileData`

- `CustomClasses/ExternalData/*` → `internal/externaldata`
  - `CountryDictionary` → `externaldata.CountryDictionary`
  - `ExtCargo` → `externaldata.ExtCargo`
  - `ExtCompany` → `externaldata.ExtCompany`
  - `LevelNames` → `externaldata.LevelNames`
  - `Routes` → `externaldata.Routes`
  - `ScsFont` → `externaldata.ScsFont`

- `CustomClasses/Global/*` → `internal/global`
  - `EnumerableExt` → utilitaires à venir dans `global` (helpers sur les slices, etc.)
  - `FloatList` → type utilitaire à venir dans `global`

- `CustomClasses/Utilities/*` → `internal/util`
  - `NumericUtilities` → `util.IntegerToHexString` et futurs helpers numériques
  - `IO_Utilities` → futurs helpers d’E/S dans `util`
  - `TextUtilities` → futurs helpers de chaînes dans `util`
  - `Web_Utilities` → futurs helpers HTTP dans `util`
  - `ZipDataUtilitiescs` → futurs helpers ZIP dans `util`
  - `Graphics` → documentation/logiciels non UI dans `util`

- `CustomClasses/UtilitiesExternal/*` → `internal/utilext`
  - `DDSImageParser` → `utilext.DDSImageParser` (stub)
  - `TGASharpLib` → `utilext.TGASharpLib` (stub)
  - `FlexibleMessageBox` → `utilext.FlexibleMessageBox` (stub)

Pour les nouveaux développements en Go, il est conseillé de :

- Utiliser les types `internal/save/...` et `internal/externaldata` plutôt que les
  classes C# originales.
- Ajouter les champs et méthodes nécessaires aux structs Go au fur et à mesure que
  la logique SII est portée (parser, sérialiseur, déchiffrement, etc.).


