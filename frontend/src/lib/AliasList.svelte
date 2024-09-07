<script lang="ts">
	import { createEventDispatcher } from "svelte";
	import { baseUrl } from "../config";
	import type { AliasResponse } from "../types";

	const aliasesUrl = baseUrl + "/v1/aliases";
	const dispatch = createEventDispatcher();

	export let aliases: AliasResponse[] = [];
	let isDeleting = false;

	function isAliasInList(alias: string) {
		return aliases.some((a) => a.alias === alias);
	}

	async function removeAlias(alias: string) {
		if (!isAliasInList(alias)) {
			return;
		}

		// TODO: Confirm dialog
		isDeleting = true;

		try {
			const response = await fetch(
				aliasesUrl + "/" + encodeURIComponent(alias),
				{
					method: "DELETE",
				},
			);

			if (response.status === 204) {
				dispatch("refresh");
			} else {
				console.error(response);
			}
		} catch (error) {
			console.error(error);
		}

		isDeleting = false;
	}
</script>

<table>
	<thead>
		<tr>
			<th scope="col">Alias</th>
			<th scope="col">E-mail</th>
			<th scope="col"></th>
		</tr>
	</thead>
	<tbody>
		{#each aliases as { alias, email }}
			<tr>
				<td>{alias}</td>
				<td>{email}</td>
				<td>
					<button disabled={isDeleting} on:click={() => removeAlias(alias)}>
						Delete
					</button>
				</td>
			</tr>
		{/each}
	</tbody>
</table>

<style></style>
