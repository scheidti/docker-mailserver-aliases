<script lang="ts">
	import { createEventDispatcher } from "svelte";

	const dispatch = createEventDispatcher();

	export let title = "";
	export let description = "";
	export let open = false;
	let modal: HTMLDialogElement;

	$: {
		if (open) {
			modal?.showModal();
		} else {
			modal?.close();
		}
	}

	function handleSubmit(event: SubmitEvent) {
		const button = event.submitter as HTMLButtonElement;

		if (button.textContent === "Confirm") {
			dispatch("confirm");
		}
		open = false;
	}
</script>

<dialog bind:this={modal} class="modal">
	<div class="modal-box">
		<h3 class="text-lg font-bold">{title}</h3>
		<p class="py-4">{description}</p>
		<div class="modal-action">
			<form method="dialog" on:submit|preventDefault={handleSubmit}>
				<button class="btn">Cancel</button>
				<button class="btn btn-primary">Confirm</button>
			</form>
		</div>
	</div>
</dialog>
