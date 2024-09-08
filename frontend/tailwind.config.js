/** @type {import('tailwindcss').Config} */
import daisyui from "daisyui"
import { emerald } from "daisyui/src/theming/themes";

export default {
	content: ["./src/**/*.{svelte,ts}", "./index.html"],
	theme: {
		extend: {},
	},
	plugins: [
		daisyui,
	],
	daisyui: {
		themes: [
			{
        dma: {
          ...emerald,
          primary: "#04A777",
					neutral: "#291720",
					"base-content": "#291720",
					secondary: "#F75C03",
        },
      },
		],
	}
};
