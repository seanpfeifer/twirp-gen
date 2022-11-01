// Code generated by protoc-gen-twirpts. DO NOT EDIT.

let host: string = "";

// setHost is useful if the API you're calling is not on the same host as your web server
export function setHost(newHost: string) {
	host = newHost;
}

// TODO: These aren't very good, and are difficult to use - make this better
// To iterate through this, use eg `...Object.entries(myMap)`
export interface StringMap<V> {
  [key: string]: V;
}

export interface NumberMap<V> {
  [key: number]: V;
}

export enum TwirpErrorCode {
	NoError = "NoError",
	Canceled = "Canceled",
	Unknown = "Unknown",
	Invalid_Argument = "Invalid_Argument",
	Malformed = "Malformed",
	Deadline_Exceeded = "Deadline_Exceeded",
	Not_Found = "Not_Found",
	Bad_Route = "Bad_Route",
	Already_Exists = "Already_Exists",
	Permission_Denied = "Permission_Denied",
	Unauthenticated = "Unauthenticated",
	Resource_Exhausted = "Resource_Exhausted",
	Failed_Precondition = "Failed_Precondition",
	Aborted = "Aborted",
	Out_Of_Range = "Out_Of_Range",
	Unimplemented = "Unimplemented",
	Internal = "Internal",
	Unavailable = "Unavailable",
	Data_Loss = "Data_Loss",
}

export class TwirpError extends Error {
	constructor(
		public code: TwirpErrorCode,
		public message: string,
		public meta?: any,
	) {
		super(code + ": " + message);
	}
}

// Patch BigInt support into JSON.stringify
BigInt.prototype["toJSON"] = function () { return this.toString() }

// Use a custom parser to allow for bigints to be properly created
async function parseJSON(res: Response): Promise<any> {
	const text = await res.text();
	return JSON.parse(text, (key, value) => {
		// For anything more than 15 digits, assume it's a bigint.
		// There is an edge case where a number with 15 digits or more will be treated as a bigint.
		if (typeof value === "string" && value.match(/^\d{15,}$/)) {
			return BigInt(value);
		}
		return value;
	});
}

function createRequest(url: string, body: any): Request {
	return new Request(host + url, {
		method: "POST",
		credentials: "same-origin",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify(body),
	});
}

export async function getCharacter(args: CharacterRequest): Promise<CharacterResponse> {
	const res = await fetch(createRequest("/rpc/characters.Characters/GetCharacter", args));
	const jsonBody = await parseJSON(res);
	if (res.ok) {
		return jsonBody;
	}
	throw new TwirpError(jsonBody.code, jsonBody.msg, jsonBody.meta);
}

export interface CharacterRequest {
	name?: string;
	mapName?: string;
	wtbItem?: Item;
	dreamClass?: Class;
}

export interface CharacterResponse {
	items?: StringMap<Item>;
	stash?: Item;
	charClass?: Class;
	charType?: CharacterResponse_Type;
	history?: CharacterResponse_History;
}

export interface CharacterResponse_History {
	parent?: string;
	hometown?: string;
}

export interface Item {
	name?: string;
	quantity?: number;
	itemType?: Item_Type;
}

export interface Item_Type {
	name?: string;
	power?: number;
}

export enum CharacterResponse_Type {
	CHARTYPE_PLAYER = "CHARTYPE_PLAYER",
	CHARTYPE_NPC = "CHARTYPE_NPC",
}

export enum Class {
	CLASS_WARRIOR = "CLASS_WARRIOR",
	CLASS_MAGE = "CLASS_MAGE",
	CLASS_ROGUE = "CLASS_ROGUE",
}
