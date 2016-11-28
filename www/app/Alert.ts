/// <reference path="./References.d.ts"/>
import * as AlertTypes from './types/AlertTypes';
import * as AlertActions from './actions/AlertActions';

export function info(message: string): void {
	AlertActions.create(AlertTypes.INFO, message);
}

export function warning(message: string): void {
	AlertActions.create(AlertTypes.WARNING, message);
}

export function error(message: string): void {
	AlertActions.create(AlertTypes.ERROR, message);
}
