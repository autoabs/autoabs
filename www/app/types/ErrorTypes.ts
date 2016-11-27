/// <reference path="../References.d.ts"/>
export const CREATE = Symbol('error.create');
export const REMOVE = Symbol('error.remove');

export interface Error {
	id?: string;
	level?: Symbol;
	message?: string;
}

export type Errors = {[key: string]: Error};
