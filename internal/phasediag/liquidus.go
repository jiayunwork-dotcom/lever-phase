package phasediag

import "fmt"

// liquidLeft 返回液相线左段折线（从纯组元 A 到共晶点）。
// 左段由 Liquidus 的前两个结点组成。
func (d PhaseDiagram) liquidLeft() Polyline {
	return Polyline{nodes: d.Liquidus.nodes[:2]}
}

// liquidRight 返回液相线右段折线（从共晶点到纯组元 B）。
func (d PhaseDiagram) liquidRight() Polyline {
	n := len(d.Liquidus.nodes)
	return Polyline{nodes: d.Liquidus.nodes[n-2:]}
}

// LiquidLeftAt 求给定温度在液相线左段上的成分（二分求根）。
// 温度必须落在左段温度跨度（共晶温度到 A 熔点）内。
func (d PhaseDiagram) LiquidLeftAt(t float64) (float64, error) {
	if err := ValidateTemperature(t); err != nil {
		return 0, err
	}
	return d.liquidLeft().CompositionAt(t, 0)
}

// LiquidRightAt 求给定温度在液相线右段上的成分（二分求根）。
func (d PhaseDiagram) LiquidRightAt(t float64) (float64, error) {
	if err := ValidateTemperature(t); err != nil {
		return 0, err
	}
	return d.liquidRight().CompositionAt(t, 0)
}

// LiquidusTempsAt 求给定温度下液相线在左右两段上的全部成分解，
// 返回的成分升序排列。温度低于共晶温度时没有解。
func (d PhaseDiagram) LiquidusTempsAt(t float64) ([]float64, error) {
	if err := ValidateTemperature(t); err != nil {
		return nil, err
	}
	return d.Liquidus.CompositionAtAll(t)
}

// LiquidusLeftTemperature 返回液相线左段温度跨度。
func (d PhaseDiagram) LiquidusLeftTemperature() (float64, float64) {
	left := d.liquidLeft()
	return left.MinTemperature(), left.MaxTemperature()
}

// LiquidusRightTemperature 返回液相线右段温度跨度。
func (d PhaseDiagram) LiquidusRightTemperature() (float64, float64) {
	right := d.liquidRight()
	return right.MinTemperature(), right.MaxTemperature()
}

// String 输出相图关键参数的描述。
func (d PhaseDiagram) String() string {
	return fmt.Sprintf("Pb-Sn phase diagram: TA=%.3fK TB=%.3fK TE=%.3fK cE=%.3f alphaMax=%.3f betaMin=%.3f",
		d.Liquidus.NodeAt(0).Temperature, d.Liquidus.NodeAt(2).Temperature, d.EutecticT, d.EutecticC, d.AlphaMax, d.BetaMin)
}
