/// <reference path="../References.d.ts"/>
export const LOADING = Symbol('build.loading');
export const LOADED = Symbol('build.loaded');
export const SYNC = Symbol('build.sync');
export const REMOVE = Symbol('build.remove');

export interface Build {
	id: string;
	name?: string;
	builder?: string;
	start?: string;
	stop?: string;
	state?: string;
	version?: string;
	release?: string;
	repo?: string;
	arch?: string;
	log?: string[];
}

export type Builds = Build[];

export interface BuildDispatch {
	type: Symbol;
	data?: {
		id?: string;
		content?: string;
		builds?: Build[];
	};
}
