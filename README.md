# a-dns

CLI d'administration des zones DNS OVH, optimisé pour une utilisation par agent IA avec sortie JSON structurée.

## Installation

```bash
go install
```

## Configuration

### Methode 1: Application Keys (traditionnel)

1. Créez une application OVH sur https://eu.api.ovh.com/createApp
2. Notez l'Application Key et le Application Secret
3. Créez un Consumer Key avec les droits DNS nécessaires
4. Configurez le CLI:

```bash
./a-dns setup
```

### Methode 2: OAuth2 Service Account (recommandé)

1. Dans OVH Manager → Gestion des comptes → Service Accounts
2. Créez un compte service avec permissions DNS (IAM)
3. Notez le Client ID et Client Secret
4. Configurez le CLI:

```bash
./a-dns setup -m 2
```

### Configuration manuelle

Créez `~/.a-dns.yaml`:

**Application Keys:**
```yaml
endpoint: ovh-eu
app_key: YOUR_APP_KEY
app_secret: YOUR_APP_SECRET
consumer_key: YOUR_CONSUMER_KEY
default_zone: ganima.xyz
```

**OAuth2 Service Account:**
```yaml
endpoint: ovh-eu
oauth2_client_id: YOUR_CLIENT_ID
oauth2_client_secret: YOUR_CLIENT_SECRET
default_zone: ganima.xyz
```

## Utilisation

### Liste des zones DNS

```bash
./a-dns list-zones
```

### Liste des enregistrements d'une zone

```bash
./a-dns list-records ganima.xyz
./a-dns list-records -t A ganima.xyz
```

### Ajouter un enregistrement

```bash
./a-dns add-record ganima.xyz A 1.2.3.4
./a-dns add-record ganima.xyz CNAME target.example.com
./a-dns add-record -s www ganima.xyz A 1.2.3.4
```

### Mettre à jour un enregistrement

```bash
./a-dns update-record ganima.xyz 123456 A 5.6.7.8
```

### Supprimer un enregistrement

```bash
./a-dns delete-record ganima.xyz 123456
```

## Format de sortie

Par défaut, la sortie est en JSON. Utilisez le flag `-o`:

```bash
./a-dns list-records ganima.xyz -o json
./a-dns list-records ganima.xyz -o yaml
```

## Exemple d'utilisation par agent IA

```
IA: Liste les enregistrements DNS de ganima.xyz
-> ./a-dns list-records ganima.xyz
-> {"zone": "ganima.xyz", "records": [...]}

IA: Ajoute un record A pour test.ganima.xyz pointant vers 1.2.3.4
-> ./a-dns add-record -s test ganima.xyz A 1.2.3.4
-> {"id": 123456, "zone": "ganima.xyz", "subDomain": "test", "type": "A", "target": "1.2.3.4", "ttl": 0}
```
