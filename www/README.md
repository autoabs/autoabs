### autoabs-www

Requires [jspm](https://www.npmjs.com/package/jspm)

```
npm install
jspm install
sed -i 's|lib/node/index.js|lib/client.js|g' jspm_packages/npm/superagent@*.js
tsc --watch
```

#### lint

```
tslint -c tslint.json app/**/*.ts*
```

#### clean

```
find app/ -name "*.js*" -delete
```

### development

```
jspm depcache app/App.js
```

#### production

```
tsc
jspm bundle app/App.js
mkdir -p dist/static
cp node_modules/@blueprintjs/core/dist/blueprint.css dist/static/
cp node_modules/material-design-icons/iconfont/material-icons.css dist/static/
cp node_modules/material-design-icons/iconfont/MaterialIcons-Regular.eot dist/static/
cp node_modules/material-design-icons/iconfont/MaterialIcons-Regular.ijmap dist/static/
cp node_modules/material-design-icons/iconfont/MaterialIcons-Regular.svg dist/static/
cp node_modules/material-design-icons/iconfont/MaterialIcons-Regular.ttf dist/static/
cp node_modules/material-design-icons/iconfont/MaterialIcons-Regular.woff dist/static/
cp node_modules/material-design-icons/iconfont/MaterialIcons-Regular.woff dist/static/
cp node_modules/material-design-icons/iconfont/MaterialIcons-Regular.woff2 dist/static/
cp styles/global.css dist/static/
cp node_modules/@blueprintjs/core/dist/blueprint.css dist/static/
cp jspm_packages/system.js dist/static/
mv build.js dist/static/app.js
mv build.js.map dist/static/app.js.map
```

#### intellij settings

Languages & Frameworks -> TypeScript: `Use TypeScript Service`

Languages & Frameworks -> TypeScript: `Uncheck Track changes`

Languages & Frameworks -> TypeScript: `Use tsconfig.json`

Usage Scope: `!file:*.js`
