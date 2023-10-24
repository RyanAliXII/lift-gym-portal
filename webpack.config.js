const path = require("path");
const Webpack = require("webpack");
const fsw = require("@nodelib/fs.walk");

const getEntries = (filesPath) => {
  const files = fsw.walkSync(filesPath, { pathSegmentSeparator: "/" });
  const entries = {};
  files.forEach((entry) => {
    const isDirectory = entry.dirent.isDirectory();
    if (!isDirectory) {
      const filename = entry.dirent.name;
      if (path.extname(filename) === ".js") {
        const newPath = `.${entry.path.replace(filesPath, "")}`;
        entries[newPath] = entry.path;
      }
    }
  });
  return entries;
};

module.exports = {
  entry: getEntries("./views"),
  mode: "development",
  output: {
    filename: "[name]",
    path: path.resolve(__dirname, "assets/js"),
  },
  resolve: {
    alias: {
      vue: "vue/dist/vue.esm-bundler.js",
    },
  },
  plugins: [
    new Webpack.DefinePlugin({
      __VUE_OPTIONS_API__: true,
      __VUE_PROD_DEVTOOLS__: true,
    }),
  ],
  module: {
    rules: [
      {
        test: /\.css$/i,
        use: ["style-loader", "css-loader"],
      },
    ],
  },
};
