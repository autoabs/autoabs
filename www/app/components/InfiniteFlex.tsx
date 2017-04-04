/// <reference path="../References.d.ts"/>
import * as React from 'react';

type BuildItem = (index: number, item: any) => JSX.Element;

type Traverse = (index: number) => void;

interface Props {
	style: React.CSSProperties;
	width: number;
	height: number;
	padding: number;
	scrollMargin: number;
	scrollMarginHit: number;
	buildItem: BuildItem;
	traverse: Traverse;
	items: any[];
	index: number;
	count: number;
}

export default class InfiniteFlex extends React.Component<Props, null> {
	ready: boolean;
	upper: number;
	upperHit: number;
	lower: number;
	lowerHit: number;
	columns: number;
	shown: number;
	index: number;

	constructor(props: any, context: any) {
		super(props, context);
		this.ready = false;
		this.upper = 0;
		this.upperHit = 0;
		this.lower = 0;
		this.lowerHit = 0;
		this.columns = 0;
		this.shown = 0;
		this.index = 0;
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
		let scroll = Math.max(window.scrollY, 0);
		let inner = window.innerHeight;
		let height = document.body.scrollHeight;
		let pos = 0;
		if (inner !== height) {
			pos = Math.max((scroll / (height - inner)) || 0, 0);
		}
		let count = this.props.count || 0;

		let elem = this.refs.container as Element;
		let width = parseInt(
			window.getComputedStyle(elem).width, 10) - this.props.padding * 2;
		this.columns = Math.floor(width / this.props.width);

		this.shown = Math.ceil(
				window.innerHeight / this.props.height) * this.columns;

		this.upper = Math.floor(
			count * pos - this.shown * this.props.scrollMargin);
		this.lower = Math.floor(
			count * pos + this.shown * this.props.scrollMargin);
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
			if (upper > 0) {
				upper -= upper % this.columns;
			}
			let lower = this.lower;
			lower += this.columns - ((lower - upper) % this.columns);
			lower = Math.min(this.props.count || 0, lower);

			style = {
				...this.props.style,
				paddingTop: ((Math.floor(upper / this.columns) * this.props.height) +
					this.props.padding) + 'px',
				paddingBottom: ((Math.floor(
					(this.props.count - lower) / this.columns) * this.props.height) +
					this.props.padding) + 'px',
			};

			let index = 0;
			for (let i = upper; i < lower; i++) {
				let item = items[i - this.props.index];

				if (item) {
					itemsDom.push(this.props.buildItem(index, item));
				}

				index += 1;
			}

			let len = items.length;
			if (len) {
				let start = this.index;
				let end = start + len;

				if ((this.props.index !== 0 && upper - start < 50)) {
					let newIndex = Math.max(0, upper - Math.floor(len / 2));
					if (newIndex !== this.index) {
						this.index = newIndex;
						setTimeout(() => {
							this.props.traverse(newIndex);
						});
					}
				} else if (end < this.props.count && end - lower < 50) {
					let newIndex = lower - Math.floor(len / 2);
					if (newIndex !== this.index) {
						this.index = newIndex;
						setTimeout(() => {
							this.props.traverse(newIndex);
						});
					}
				}
			}
		}

		return <div ref="container" style={style}>
			{itemsDom}
		</div>;
	}
}
