export type AliasResponse = {
	alias: string;
	email: string;
};

export type AliasListResponse = {
	aliases: AliasResponse[];
};

export type EmailsListResponse = {
	emails: string[];
};

export type ErrorResponse = {
	error: string;
};

export type StatusResponse = {
	running: boolean;
};

export type Toast = {
	text: string;
	type: "error" | "success" | "info" | "warning";
};
