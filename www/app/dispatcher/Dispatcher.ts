/// <reference path="../References.d.ts"/>
import * as Flux from 'flux';
import * as GlobalTypes from '../types/GlobalTypes';

class _Dispatcher extends Flux.Dispatcher<GlobalTypes.Dispatch> {}
let Dispatcher = new _Dispatcher();
export default Dispatcher;
