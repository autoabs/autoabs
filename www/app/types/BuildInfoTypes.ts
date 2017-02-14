/// <reference path="../References.d.ts"/>
export const OPEN = Symbol('build_info.open');
export const CLOSE = Symbol('build_info.close');
export const UPDATE = Symbol('build_info.update');

export interface BuildInfoDispatch {
	type: Symbol;
	data?: {
		id?: string;
		output?: string[];
	};
}
