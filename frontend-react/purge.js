const purify = require("purify-css")

var content = ['**/src/routes/*.jsx'];
var css = ['**/src/css/*.css'];

const opts = {
    output: 'purified.css'
};


purify(content, css, opts, function (res) {
    log(res);
});