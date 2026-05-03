# ZINTENT + Laravel (Blade)

Use this pattern in a full Laravel app: compile CSS from Blade views + any front-end sources you keep under `resources/`.

## Idea

1. The ZINTENT compiler reads **one entry** HTML/Blade file (or a small shell) and scans **`--content`** for all classes.
2. Output a single `public/css/zintent.css` (or `resources/css/zintent.css` if you prefer).
3. Reference that file from your main layout Blade.

## Example commands (run from ZINTENT repo root)

Adjust paths to match your Laravel project if you copy this layout:

```bash
go run ./compiler ^
  --input examples/laravel-blade/resources/views/welcome.blade.php ^
  -o examples/laravel-blade/public/css/zintent.css ^
  --preset core ^
  --scanner parser ^
  --content examples/laravel-blade/resources/views
```

On macOS/Linux use `\` line continuation instead of `^`.

## Laravel integration

- In `resources/views/layouts/app.blade.php` (or similar):

```blade
<link rel="stylesheet" href="{{ asset('css/zintent.css') }}">
```

- Add an npm/composer script or CI step that runs the ZINTENT `go run ./compiler` command whenever Blade or JS views change.

## Notes

- `--scanner parser` helps for Blade files that are mostly HTML.
- If you use Vue/React inside Laravel (e.g. Inertia), add those source folders to `--content`.
- See `docs/MIGRATION_FROM_TAILWIND.md` if migrating utility classes.
