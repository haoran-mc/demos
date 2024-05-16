<!-- # React + Vite

This template provides a minimal setup to get React working in Vite with HMR and some ESLint rules.

Currently, two official plugins are available:

- [@vitejs/plugin-react](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react/README.md) uses [Babel](https://babeljs.io/) for Fast Refresh
- [@vitejs/plugin-react-swc](https://github.com/vitejs/vite-plugin-react-swc) uses [SWC](https://swc.rs/) for Fast Refresh -->



```shell
# 使用 yarn 管理包
npm install --global yarn

# 创建一个 vite 项目
yarn create vite
```

```
SWC is a free and open-source JavaScript transcompiler like Babel, but unlike Babel it can significantly speed up build and development time due to SWC’s fast conversion capabilities, however it may not support all Babel plugins.

SWC plays an integral role in the Vite ecosystem. Vite uses SWC to run the Babel transformation pipeline much faster, which is especially useful for large React projects.
```

但是我并不打算在一个 demo 里使用 SWC。

```
❯ yarn create vite

...
success Installed "create-vite@5.2.3" with binaries:
      - create-vite
      - cva
✔ Project name: … react-vite-tailwind
✔ Select a framework: › React
✔ Select a variant: › JavaScript
...
```

然后运行 `yarn dev` 运行测试。


```shell
# 下载 tailwindcss 要求的依赖
yarn add -D tailwindcss postcss autoprefixer

# 初始化 tailwindcss
npx tailwindcss init -p
# 会创建两个配置文件
# 1. tailwind.config.js tailwind 的配置文件
# 2. postcss.config.js postcss 的配置文件
```

-----

```
└── src/
    ├── pages/           # 页面组件，例如首页、关于页面、联系页面
    ├── widgets/         # 解决方案部分、解决流程部分
    └── component/       # 页面共用
```

`yarn dev` 运行。