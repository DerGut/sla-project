const path = require('path');

module.exports = {
    mode: 'development',
    entry: './static/lib/index.js',
    output: {
        filename: 'main.js',
        path: path.resolve(__dirname, 'static/js'),
    },
    externals: {
        'react': 'React'
    }
};
