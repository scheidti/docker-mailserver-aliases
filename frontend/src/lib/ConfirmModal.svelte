<script lang="ts">
	import { onMount } from "svelte";

	interface Props {
		title?: string;
		description?: string;
		open?: boolean;
		confirm?: () => void;
	}

	let { title = "", description = "", open = $bindable(false), confirm }: Props = $props();
	let modal: HTMLDialogElement | undefined = $state();

	$effect(() => {
		if (open) {
			modal?.showModal();
		} else {
			modal?.close();
		}
	});

	function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		const button = event.submitter as HTMLButtonElement;
		if (button.textContent === "Confirm") {
			confirm?.();
		}
		open = false;
	}

	onMount(() => {
		modal?.addEventListener("close", (event) => {
			event.preventDefault();
			open = false;
		});
	});
</script>

<dialog bind:this={modal} class="modal" title={title}>
	<div class="modal-box">
		<h3 class="text-lg font-bold">{title}</h3>
		<p class="py-4">{description}</p>
		<div class="modal-action">
			<form method="dialog" onsubmit={handleSubmit}>
				<button class="btn">Cancel</button>
				<button class="btn btn-primary">Confirm</button>
			</form>
		</div>
	</div>
</dialog>
