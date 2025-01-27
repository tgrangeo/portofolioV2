# La communication entre htmx (client) et Go (serveur) : Les bases

Dans le développement web moderne, **htmx** s'impose comme un outil léger et efficace pour créer des interfaces utilisateur dynamiques sans avoir besoin d'utiliser un framework JavaScript lourd. Associé à un backend Go, il permet de créer des applications performantes et maintenables. Cet article explique les bases de la communication entre htmx (côté client) et Go (côté serveur) avec des exemples concrets.

---

## Qu'est-ce que htmx ?

**htmx** est une bibliothèque JavaScript minimaliste qui permet d'effectuer des requêtes HTTP directement à partir du HTML, en exploitant des attributs tels que `hx-get`, `hx-post`, `hx-swap`, etc. Cela simplifie la création de fonctionnalités dynamiques sans avoir à écrire de JavaScript complexe.

## Exemple de scénario

Imaginons une application simple où l'utilisateur peut ajouter et afficher des tâches dynamiquement via un formulaire. Le backend Go gère les requêtes, tandis que htmx s'occupe des interactions côté client.

---

## Mise en place

### 1. **Structure HTML de base**

Voici une page HTML minimaliste intégrant htmx pour interagir avec le serveur Go :

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Gestion des tâches</title>
    <script src="https://unpkg.com/htmx.org"></script>
</head>
<body>
    <h1>Liste des tâches</h1>
    <div id="task-list">
        <!-- Les tâches seront chargées ici -->
    </div>

    <h2>Ajouter une tâche</h2>
    <form hx-post="/tasks" hx-target="#task-list" hx-swap="beforeend">
        <input type="text" name="task" placeholder="Nouvelle tâche" required>
        <button type="submit">Ajouter</button>
    </form>
</body>
</html>
