/// <reference path="../References.d.ts"/>
export const CHANGE = Symbol('change');

export interface Dispatch {
	type: Symbol;
	data?: any;
}
