const path = require('path');

module.exports = {
    entry: './static/script.js', // Входной файл
    output: {
        path: path.resolve(__dirname, 'static'), // Папка для сборки
        filename: 'bundle.js', // Имя выходного файла
    },
    mode: 'production', // Минимизация и оптимизация
};