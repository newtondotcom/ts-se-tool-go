- change the structure of the sii* files and folder


- pour la cli, utiliser https://github.com/urfave/cli 
La premiere entrée attendue est un message de bienvenue 
Ensuite détection des jeux installés, si ATS et ETS2 trouvés, demander au joueur de choisir parmi quel jeu puis quel profil
si 0 jeu trouvé demander à l'utilisateur le path de son profil directement
Sauvegarde du profil `mv slot slot.bckp`
Ensuite, voici les features à implémenter : 
+ mettre somme argent
+ mettre points xp
+ mettre tous les skills au max
+ acheter tous les garages 
+ améliorer tous les garages
+ peupler tous les garages avec des camions random
+ recruter des employés et peupler tous les camions


- - backup the current profile
- Manage saves and put 1/game.sii and 1/info.sii
```md
Analyse de l'écriture des fichiers `.sii` (game.sii, info.sii, profile.sii) dans TS SE Tool.

## Ce que fait TS SE Tool lors de l'écriture des fichiers `.sii`

### 1. **`game.sii`** — Sauvegarde principale

**Fichier** : `MethodsReadWrite.cs` → `NewWrireSaveFile()`

**Processus** :

- **Chemin déterminé** :
  ```csharp
  SavefilePath = Globals.SavesHex[comboBoxSaves.SelectedIndex];
  string SiiSavePath = SavefilePath + @"\game.sii";
  ```
  - `SavefilePath` vient de la liste des sauvegardes détectées dans le profil sélectionné.
  - Exemple : `C:\Users\...\Documents\Euro Truck Simulator 2\profiles\4D792050726F66696C65\save\1\game.sii`

- **Nom du dossier** : Le dossier de sauvegarde est celui du jeu (ex: `save\1`, `save\autosave`, etc.). TS SE Tool ne renomme pas ces dossiers.

- **Écriture** :
  ```csharp
  using (StreamWriter writer = new StreamWriter(SiiSavePath, false))
  {
      writer.Write(SiiNunitData.PrintOut(0));
  }
  ```
  - Le `false` écrase le fichier existant.
  - `SiiNunitData.PrintOut(0)` régénère tout le contenu SiiNunit depuis la structure en mémoire.

- **Suppression de dossiers** : Non. Seul le fichier `game.sii` est écrasé dans le dossier existant.

---

### 2. **`info.sii`** — Métadonnées de sauvegarde

**Fichier** : `MethodsReadWrite.cs` (lecture) / `FormConvoyControlPositions.cs` (écriture)

**Processus** :

- **Écriture normale** : Dans le même dossier que `game.sii` :
  ```csharp
  string SiiInfoPath = SavefilePath + @"\info.sii";
  ```

- **Écriture pour convoy tools** (création de nouvelles sauvegardes) :
  ```csharp
  string fp = Directory.GetParent(Globals.SavesHex[0]).FullName + "\\" + NewCustomFolders[iSave].ToString();
  Directory.CreateDirectory(fp);  // Crée un NOUVEAU dossier
  using (StreamWriter writer = new StreamWriter(fp + "\\info.sii", false))
  {
      infoData.WriteToStream(writer);
  }
  ```
  - Un nouveau dossier de sauvegarde est créé (nommé selon `NewCustomFolders`).
  - `info.sii` est écrit dedans.

- **Suppression de dossiers** : Non, sauf si tu supprimes manuellement une sauvegarde via l’outil convoy.

---

### 3. **`profile.sii`** — Données du profil

**Fichier** : `Forms/ProfileEditor/FormProfileEditorRenameClone.cs`

**Processus** :

- **Chemin déterminé** :
  ```csharp
  string SiiProfilePath = Globals.ProfilesHex[comboBoxProfiles.SelectedIndex] + @"\profile.sii";
  ```
  - Le dossier de profil est celui du jeu (nom en hex, ex: `4D792050726F66696C65`).

- **Renommage de profil** :
  1. Création d’un nouveau dossier (nom hex du nouveau nom) :
     ```csharp
     NewFolderName = Utilities.TextUtilities.FromStringToHex(NewProfileName);
     NewFolderPath = InitialPath.Remove(...) + NewFolderName;
     Directory.CreateDirectory(NewFolderPath);
     ```
  2. Copie complète du contenu :
     ```csharp
     Utilities.IO_Utilities.DirectoryCopy(InitialPath, NewFolderPath, true);
     ```
  3. Écriture de `profile.sii` modifié :
     ```csharp
     using (StreamWriter SW = new StreamWriter(NewFolderPath + "\\profile.sii", false))
     {
         ProfileData.WriteToStream(SW);
     }
     ```
  4. Suppression de l’ancien dossier :
     ```csharp
     Directory.Delete(InitialPath, true);
     ```

- **Suppression de dossiers** : Oui, uniquement lors du renommage (l’ancien dossier est supprimé après succès).

---

### Résumé

| Fichier | Dossier utilisé | Création dossier | Suppression dossier |
|---------|------------------|------------------|---------------------|
| **`game.sii`** | Dossier de sauvegarde existant (ex: `save\1`) | Non | Non |
| **`info.sii`** | Même dossier que `game.sii` (ou nouveau si convoy) | Seulement si création nouvelle sauvegarde | Non |
| **`profile.sii`** | Dossier de profil existant (ex: `4D792050726F66696C65`) | Seulement si renommage/clonage | Seulement si renommage (ancien supprimé) |

**Points importants** :
- TS SE Tool n’invente pas de nouveaux noms de dossiers : il utilise ceux du jeu (profils en hex, saves numérotés).
- Il écrase les fichiers `.sii` existants (mode `false` sur `StreamWriter`).
- Il ne supprime pas de dossiers, sauf lors du renommage de profil (ancien supprimé après copie réussie).
- Les fichiers sont écrits en texte SiiNunit (non chiffré) après décodage.
```