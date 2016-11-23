/// <reference path="../References.d.ts"/>
export function uuid(): string {
	return (+new Date() + Math.floor(Math.random() * 999999)).toString(36);
}
