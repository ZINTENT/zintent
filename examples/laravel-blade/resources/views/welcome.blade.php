<!DOCTYPE html>
<html lang="{{ str_replace('_', '-', app()->getLocale()) }}">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>ZINTENT Laravel example</title>
    @if(file_exists(public_path('css/zintent.css')))
        <link rel="stylesheet" href="{{ asset('css/zintent.css') }}">
    @endif
</head>
<body class="zi-reset zi-base-body">
    <main class="zi-container zi-py-12 intent-stack-md">
        <h1 class="zi-text-2xl zi-font-bold">ZINTENT + Blade</h1>
        <p class="zi-text-muted">Compile CSS with the <code class="zi-p-2">go run ./compiler</code> command from the README in this folder.</p>
    </main>
</body>
</html>
