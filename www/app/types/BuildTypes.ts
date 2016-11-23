/// <reference path="../References.d.ts"/>
export const LOADING = Symbol('build.loading');
export const LOAD = Symbol('build.load');
export const CREATE = Symbol('build.create');
export const REMOVE = Symbol('build.remove');
export const UPDATE = Symbol('build.update');

export interface Build {
	id: string;
	name?: string;
	builder?: string;
	start?: Date;
	stop?: Date;
	state?: string;
	version?: string;
	release?: string;
	repo?: string;
	arch?: string;
	log?: string[];
}

export type Builds = {[key: string]: Build}

export interface BuildDispatch {
	type: Symbol;
	data?: {
		id: string;
		content: string;
		builds: Build[];
	};
}
