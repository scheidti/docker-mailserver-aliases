<script lang="ts">
	import AddAlias from "./lib/AddAlias.svelte";
	import Alert from "./lib/Alert.svelte";
	import AliasList from "./lib/AliasList.svelte";
	import Spinner from "./lib/Spinner.svelte";
	import { baseUrl } from "./config";
	import type {
		AliasListResponse,
		AliasResponse,
		StatusResponse,
	} from "./types";

	const running = checkIfMailserverIsRunning();
	const aliasesUrl = baseUrl + "/v1/aliases";
	const statusUrl = baseUrl + "/v1/status";

	let isLoading = false;
	let aliases: AliasResponse[] = [];

	async function checkIfMailserverIsRunning() {
		try {
			const response = await fetch(statusUrl);
			const data: StatusResponse = await response.json();
			return data.running === true;
		} catch {
			return false;
		}
	}

	async function getAliases() {
		isLoading = true;

		try {
			const response = await fetch(aliasesUrl);
			const list: AliasListResponse = await response.json();
			aliases = list.aliases;
		} catch {}

		isLoading = false;
	}
</script>

<header>
	<h1>Docker Mailserver Aliases</h1>
</header>

<main>
	{#await running}
		<Spinner />
	{:then isRunning}
		{#if isRunning}
			<AddAlias on:added={getAliases} />
			{#if isLoading}
				<Spinner />
			{:else}
				<AliasList on:refresh={getAliases} {aliases} />
			{/if}
		{:else}
			<Alert message={"Mailserver is not running."} type={"error"} />
		{/if}
	{/await}
</main>

<footer>
	© 2024
	<a
		href="https://github.com/scheidti/docker-mailserver-aliases"
		target="_blank"
	>
		Christian Scheid
	</a>
	- This project is licensed under the MIT license.
</footer>

<style>
</style>
