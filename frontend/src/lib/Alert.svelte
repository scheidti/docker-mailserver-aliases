<script lang="ts">
	export let message: string = "";
	export let type: "error" | "success" | "info" | "warning" = "info";
</script>

<div role="alert" class={`alert alert-${type}`}>
	<span>{message}</span>
</div>

<style></style>
