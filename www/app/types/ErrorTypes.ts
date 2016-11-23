/// <reference path="../References.d.ts"/>
export const CREATE = Symbol('error_create');
export const REMOVE = Symbol('error_remove');

export interface Error {
	id?: string;
	level?: Symbol;
	message?: string;
}

export type Errors = {[key: string]: Error}
