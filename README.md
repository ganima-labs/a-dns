# a-dns

CLI d'administration des zones DNS OVH et gestion des certificats Let's Encrypt, optimisé pour une utilisation par agent IA avec sortie JSON structurée.

## Installation

```bash
go build -o a-dns .
# Ou
go install
```

## Configuration OAuth2 (Recommandé)

### Étape 1: Créer un Service Account OVH

1. Connectez-vous au [Manager OVH](https://www.ovh.com/manager/dedicated)
2. Cliquez sur votre nom en haut à droite → **Gestion des comptes**
3. Allez dans l'onglet **Comptes de service**
4. Cliquez sur **Ajouter un compte de service**
5. Remplissez:
   - **Nom**: `a-dns`
   - **Description**: `CLI DNS et certificats pour agent IA`
6. Validez la création

### Étape 2: Récupérer les credentials

⚠️ **Important**: Notez immédiatement le **Client ID** et le **Client Secret**, ils ne seront plus affichés ensuite.

Exemple:
```
Client ID: EU.a55039ff110a2cf1
Client Secret: f26125b48155160ba61fec3e3ea8912c
```

### Étape 3: Configurer les permissions IAM

1. Dans le compte de service créé, allez dans l'onglet **Politique IAM**
2. Cliquez sur **Ajouter une politique**
3. Sélectionnez **Politique personnalisée**
4. Ajoutez la politique suivante:

```json
{
  "rules": [
    {
      "effect": "allow",
      "action": "*",
      "resource": "urn:v1:eu:resource:dnsZone:*"
    }
  ]
}
```

5. Validez la politique

### Étape 4: Configurer le CLI

```bash
./a-dns setup -m 2 \
  -i EU.a55039ff110a2cf1 \
  -x f26125b48155160ba61fec3e3ea8912c \
  -z votre-domaine.com \
  -e votre@email.com
```

**Paramètres:**
- `-m 2` : Utiliser OAuth2 (Service Account)
- `-i` : Client ID
- `-x` : Client Secret
- `-z` : Zone DNS par défaut
- `-e` : Email pour Let's Encrypt

### Étape 5: Vérifier la configuration

```bash
# Lister vos zones DNS
./a-dns list-zones

# Devrait retourner quelque chose comme:
# ["votre-domaine.com", "autre-domaine.com"]
```

---

## Configuration Application Keys (Alternative)

Si vous préférez la méthode traditionnelle:

### Étape 1: Créer une application OVH

1. Allez sur https://eu.api.ovh.com/createApp
2. Connectez-vous avec votre compte OVH
3. Créez une nouvelle application:
   - **Nom**: `a-dns`
   - **Description**: `CLI DNS`
4. Notez l'**Application Key** et l'**Application Secret**

### Étape 2: Générer un Consumer Key

```bash
curl -X POST https://eu.api.ovh.com/1.0/auth/credential \
  -H 'Content-Type: application/json' \
  -H 'X-Ovh-Application: VOTRE_APP_KEY' \
  -d '{
    "accessRules": [
      {"method": "GET", "path": "/domain/zone/*"},
      {"method": "POST", "path": "/domain/zone/*"},
      {"method": "PUT", "path": "/domain/zone/*"},
      {"method": "DELETE", "path": "/domain/zone/*"}
    ]
  }'
```

### Étape 3: Valider le Consumer Key

1. Ouvrez l'URL `validationUrl` retournée
2. Connectez-vous et validez les droits
3. Notez le `consumerKey`

### Étape 4: Configurer le CLI

```bash
./a-dns setup \
  -k VOTRE_APP_KEY \
  -s VOTRE_APP_SECRET \
  -c VOTRE_CONSUMER_KEY \
  -z votre-domaine.com \
  -e votre@email.com
```

---

## Fichier de configuration

Le fichier `~/.a-dns.yaml` est créé automatiquement. Vous pouvez l'éditer manuellement:

**OAuth2 (recommandé):**
```yaml
endpoint: ovh-eu
default_zone: votre-domaine.com
email: votre@email.com
oauth2_client_id: EU.xxxxxxxxxxxxxxxx
oauth2_client_secret: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

**Application Keys:**
```yaml
endpoint: ovh-eu
default_zone: votre-domaine.com
email: votre@email.com
app_key: xxxxxxxxxxxxxxxx
app_secret: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
consumer_key: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

---

## Commandes DNS

### Lister les zones

```bash
./a-dns list-zones
```

### Lister les enregistrements

```bash
./a-dns list-records votre-domaine.com
./a-dns list-records -t A votre-domaine.com        # Filtrer par type
./a-dns list-records -s www votre-domaine.com      # Filtrer par sous-domaine
```

### Ajouter un enregistrement

```bash
./a-dns add-record votre-domaine.com A 192.168.1.1
./a-dns add-record -s www votre-domaine.com A 192.168.1.1
./a-dns add-record votre-domaine.com CNAME target.example.com
./a-dns add-record -t TXT votre-domaine.com "valeur txt"
```

### Mettre à jour un enregistrement

```bash
./a-dns update-record votre-domaine.com 123456 A 192.168.1.2
```

### Supprimer un enregistrement

```bash
./a-dns delete-record votre-domaine.com 123456
```

---

## Certificats Let's Encrypt

Le CLI gère automatiquement les certificats via challenge DNS-01 (aucun serveur web requis).

### Demander un certificat

```bash
# Certificat simple
./a-dns cert request votre-domaine.com

# Avec sous-domaines (SANs)
./a-dns cert request votre-domaine.com --sans www.votre-domaine.com,api.votre-domaine.com

# Wildcard (inclut *.votre-domaine.com et votre-domaine.com)
./a-dns cert request "*.votre-domaine.com" --sans votre-domaine.com
```

### Test avec staging (recommandé avant production)

```bash
# Utilise l'environnement de test Let's Encrypt (pas de vrai certificat)
./a-dns cert request votre-domaine.com --staging
```

### Renouveler un certificat

```bash
./a-dns cert renew votre-domaine.com
```

Le renouvellement ne fait rien si le certificat a plus de 30 jours restants (configurable avec `--renew-days`).

### Lister les certificats

```bash
./a-dns cert list
```

### Vérifier le statut

```bash
./a-dns cert status votre-domaine.com
```

### Emplacement des fichiers

Les certificats sont stockés dans `~/.a-dns/certs/`:
```
~/.a-dns/certs/
├── votre-domaine.com.crt    # Certificat
├── votre-domaine.com.key    # Clé privée
├── _.votre-domaine.com.crt  # Wildcard
├── _.votre-domaine.com.key
├── .acme-user.key           # Clé compte ACME
└── .acme-user-reg.json      # Registration ACME
```

### Options certificats

| Option | Description | Défaut |
|--------|-------------|--------|
| `--cert-dir` | Répertoire de stockage | `~/.a-dns/certs` |
| `--email` | Email pour Let's Encrypt | (depuis config) |
| `--staging` | Environnement de test | `false` |
| `--renew-days` | Seuil de renouvellement | `30` |
| `--sans` | Subject Alternative Names | `[]` |

---

## Format de sortie

Par défaut JSON (optimisé pour parsing IA):

```bash
./a-dns list-records votre-domaine.com
```

```json
[
  {
    "id": 5215115250,
    "zone": "votre-domaine.com",
    "subDomain": "www",
    "type": "A",
    "target": "192.168.1.1",
    "ttl": 3600
  }
]
```

Autres formats:
```bash
./a-dns list-records votre-domaine.com -o yaml
./a-dns list-records votre-domaine.com -o table
```

---

## Utilisation par agent IA

### Exemples de prompts

**Lister les enregistrements:**
```
IA: Liste les enregistrements DNS de votre-domaine.com
-> ./a-dns list-records votre-domaine.com
```

**Ajouter un record:**
```
IA: Ajoute un record A pour test.votre-domaine.com pointant vers 1.2.3.4
-> ./a-dns add-record -s test votre-domaine.com A 1.2.3.4
```

**Vérifier certificat:**
```
IA: Vérifie si le certificat de votre-domaine.com doit être renouvelé
-> ./a-dns cert status votre-domaine.com
-> {"daysLeft": 25, "needsRenew": true}
```

**Renouveler si nécessaire:**
```
IA: Renouvelle le certificat votre-domaine.com
-> ./a-dns cert renew votre-domaine.com
```

### Script de renouvellement automatique

```bash
#!/bin/bash
# A placer dans cron (hebdomadaire)

for domain in $(./a-dns cert list | jq -r '.[].domain'); do
  status=$(./a-dns cert status "$domain")
  needsRenew=$(echo "$status" | jq -r '.needsRenew')
  
  if [ "$needsRenew" = "true" ]; then
    echo "Renouvellement de $domain..."
    ./a-dns cert renew "$domain"
  fi
done
```

---

## Dépannage

### Erreur "missing authentication information"

Vérifiez que le fichier `~/.a-dns.yaml` contient les bons credentials:
```bash
cat ~/.a-dns.yaml
```

### Erreur "could not find zone for domain"

1. Vérifiez que la zone DNS existe dans votre compte OVH
2. Vérifiez les permissions IAM du Service Account

### Erreur "unknown record ID"

Le challenge DNS n'a pas été nettoyé correctement. Vous pouvez ignorer ce warning.

### Tester la connexion API

```bash
./a-dns list-zones
```

Si cette commande retourne vos zones, la configuration est correcte.
