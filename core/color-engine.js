/**
 * ZINTENT AI-Powered Color Engine
 * Generates harmonious palettes using color theory algorithms.
 */

function generatePalette(baseHue) {
    return {
        50:  `hsl(${baseHue}, 100%, 98%)`,
        100: `hsl(${baseHue}, 100%, 94%)`,
        200: `hsl(${baseHue}, 90%, 85%)`,
        300: `hsl(${baseHue}, 80%, 75%)`,
        400: `hsl(${baseHue}, 75%, 65%)`,
        500: `hsl(${baseHue}, 80%, 55%)`, // Primary
        600: `hsl(${baseHue}, 90%, 45%)`,
        700: `hsl(${baseHue}, 95%, 35%)`,
        800: `hsl(${baseHue}, 100%, 25%)`,
        900: `hsl(${baseHue}, 100%, 15%)`,
    };
}

module.exports = { generatePalette };
