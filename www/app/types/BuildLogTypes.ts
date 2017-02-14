/// <reference path="../References.d.ts"/>
export const OPEN = Symbol('build_log.open');
export const CLOSE = Symbol('build_log.close');
export const UPDATE = Symbol('build_log.update');

export interface BuildLogDispatch {
	type: Symbol;
	data?: {
		id?: string;
		output?: string[];
	};
}
