# a-dns — CLI DNS OVH pour Agents IA

## Description

CLI d'administration des zones DNS OVH et certificats Let's Encrypt, optimisé pour une utilisation par agent IA avec sortie JSON structurée.

## Installation

```bash
# Compiler
go build -o a-dns

# Configuration initiale (OAuth2 recommandé)
./a-dns setup
```

## Commandes

### DNS

```bash
# Lister les zones
./a-dns list-zones

# Lister les enregistrements d'une zone
./a-dns list-records example.com
./a-dns list-records example.com --subdomain www --type A

# Ajouter un enregistrement
./a-dns add-record example.com A --target "192.168.1.1" --subdomain www --ttl 3600

# Mettre à jour un enregistrement
./a-dns update-record example.com 12345 --target "192.168.1.2"

# Supprimer un enregistrement
./a-dns delete-record example.com 12345
```

### Certificats Let's Encrypt

```bash
# Demander un certificat (challenge DNS-01 automatique)
./a-dns cert request example.com --email admin@example.com

# Avec SANs
./a-dns cert request example.com --email admin@example.com --sans "www.example.com,api.example.com"

# Renouveler
./a-dns cert renew example.com

# Lister les certificats locaux
./a-dns cert list

# Statut d'un certificat
./a-dns cert status example.com
```

### Skill (pour agents IA)

```bash
# Installer la skill pour un agent spécifique
./a-dns skill install kilo
./a-dns skill install claude
./a-dns skill install cursor

# Installer pour tous les agents détectés
./a-dns skill install --all
```

## Sortie JSON

Toutes les commandes retournent du JSON structuré par défaut:

```json
{
  "success": true,
  "data": [...],
  "error": null
}
```

## Configuration

Fichier `~/.a-dns.yaml`:

```yaml
auth_type: oauth2
oauth2_client_id: "xxx"
oauth2_client_secret: "xxx"
api_endpoint: "ovh-eu"
email: "admin@example.com"
default_zone: "example.com"
language: "fr"
```

## Authentification

### OAuth2 Service Account (recommandé)

1. Manager OVH → Gestion des comptes → Comptes de service
2. Créer un compte de service
3. Ajouter une politique IAM avec droits DNS
4. Noter Client ID et Client Secret

### Application Keys (legacy)

1. https://eu.api.ovh.com/createApp
2. Générer Consumer Key avec droits DNS
3. Valider l'URL de validation

## Architecture

```
a-dns/
├── main.go                 # Entry point
├── cmd/                    # Commandes Cobra
│   ├── root.go            # Commande racine
│   ├── setup.go           # Configuration interactive
│   ├── cert.go            # Certificats Let's Encrypt
│   ├── list_*.go          # Commandes DNS
│   └── skill.go           # Installation skill
├── internal/
│   ├── config/            # Configuration Viper
│   ├── ovhclient/         # Client OVH API
│   ├── certbot/           # Wrapper ACME/lego
│   └── i18n/              # Internationalisation
└── skill/
    ├── skill.go           # Contenu embarqué
    └── SKILL.md           # Ce fichier
```

## Internationalisation

Langues supportées: `fr` (défaut), `en`, `es`

```bash
./a-dns --lang en list-zones
```

## Dépendances clés

- `github.com/spf13/cobra` — CLI framework
- `github.com/spf13/viper` — Configuration
- `github.com/ovh/go-ovh` — Client API OVH
- `github.com/go-acme/lego` — Client ACME Let's Encrypt
- `github.com/nicksnyder/go-i18n/v2` — Internationalisation

## Patterns de code

### Ajouter une commande

```go
func NewMyCmd() *cobra.Command {
    var flag string
    cmd := &cobra.Command{
        Use:   "my-cmd <arg>",
        Short: i18n.T("my.cmd.short"),
        Args:  cobra.ExactArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            result := doSomething(args[0], flag)
            outputResult(result)
        },
    }
    cmd.Flags().StringVar(&flag, "flag", "", i18n.T("flag.desc"))
    return cmd
}
```

### Appel API OVH

```go
client, _ := config.GetOVHClient()
var records []DNSRecord
client.Get(&records, "/domain/zone/example.com/record")
```

### Traduction

```go
i18n.T("key.name")
i18n.TWithData("key.template", map[string]interface{}{"Var": value})
```

## Erreurs courantes

| Erreur | Solution |
|--------|----------|
| `403 Forbidden` | Vérifier les permissions IAM/OVH |
| `404 Not Found` | Zone ou enregistrement inexistant |
| `ACME registration` | Vérifier l'email et le compte Let's Encrypt |
| `DNS propagation` | Attendre quelques minutes après modification |
