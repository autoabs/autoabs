/// <reference path="../References.d.ts"/>
import * as React from 'react';

interface BuildItem {
	(index: number, item: any): JSX.Element;
}

interface Props {
	style: React.CSSProperties;
	width: number;
	height: number;
	padding: number;
	scrollMargin: number;
	scrollMarginHit: number;
	buildItem: BuildItem;
	items: any[];
}

export default class InfiniteFlex extends React.Component<Props, null> {
	ready: boolean;
	upper: number;
	upperHit: number;
	lower: number;
	lowerHit: number;
	columns: number;
	shown: number;

	constructor(props: any, context: any) {
		super(props, context);
		this.ready = false;
		this.upper = 0;
		this.upperHit = 0;
		this.lower = 0;
		this.lowerHit = 0;
		this.columns = 0;
		this.shown = 0;
	}

	componentDidMount(): void {
		window.addEventListener('scroll', this.onScroll);
		window.addEventListener('resize', this.onScroll);
		this.ready = true;
		this.forceUpdate();
	}

	componentWillUnmount(): void {
		window.removeEventListener('scroll', this.onScroll);
		window.removeEventListener('resize', this.onScroll);
		this.ready = false;
	}

	updateScroll = (): void => {
		let scroll = window.scrollY;
		let inner = window.innerHeight;
		let height = document.body.scrollHeight;
		let pos = (scroll / (height - inner)) || 0;
		let len = this.props.items.length;

		let elem = this.refs['container'] as Element;
		let width = parseInt(
			window.getComputedStyle(elem).width) - this.props.padding * 2;
		this.columns = Math.floor(width / this.props.width);

		this.shown = Math.ceil(
				window.innerHeight / this.props.height) * this.columns;

		this.upper = Math.floor(
			len * pos - this.shown * this.props.scrollMargin);
		this.lower = Math.floor(
			len * pos + this.shown * this.props.scrollMargin);
	}

	updateScrollHit = (): void => {
		this.upperHit = this.upper - this.shown * this.props.scrollMarginHit;
		this.lowerHit = this.lower + this.shown * this.props.scrollMarginHit;
	}

	onScroll = (): void => {
		this.updateScroll();

		if (this.upper <= this.upperHit || this.lower >= this.lowerHit) {
			this.upperHit = this.upper - this.shown;
			this.lowerHit = this.lower + this.shown;
			this.forceUpdate();
		}
	}

	render(): JSX.Element {
		let style = {};
		let itemsDom: JSX.Element[] = [];

		if (this.ready) {
			this.updateScroll();
			this.updateScrollHit();
			let items = this.props.items;
			let upper = Math.max(0, this.upper);
			let lower = this.lower;
			if (upper > 0) {
				upper -= upper % this.columns;
			}
			lower += this.columns - ((lower - upper) % this.columns);
			lower = Math.min(items.length, lower);

			style = {
				...this.props.style,
				paddingTop: ((Math.floor(upper / this.columns) * 123) +
					this.props.padding) + 'px',
				paddingBottom: ((Math.floor(
					(items.length - lower) / this.columns) * 123) +
					this.props.padding) + 'px',
			};

			let index = 0;
			for (let i = upper; i < lower; i++) {
				let item = items[i];
				itemsDom.push(this.props.buildItem(index, item));
				index += 1;
			}
		}

		return <div ref="container" style={style}>
			{itemsDom}
		</div>;
	}
}
