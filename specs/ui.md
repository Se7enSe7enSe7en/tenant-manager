# Icons

## Where to get icons?

- google fonts icons: https://fonts.google.com/icons?icon.size=24&icon.color=%23e3e3e3

### download the .svg then remove the fill, width and height attributes

```html
<!-- from -->
<svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="#e3e3e3"><path ...></svg>

<!-- to -->
 <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 -960 960 960"><path ... ></svg>
```

### we remove these attributes so that we can use tailwind to override them

```html
<svg class="size-12 fill-amber-500">
  <use href="/path/to/icon.svg"></use>
</svg>
```
