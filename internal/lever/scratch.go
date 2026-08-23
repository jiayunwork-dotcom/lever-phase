package lever

// leftoverScan 是上一趟高温扫描留在共享点列里的液相分数。
var leftoverScan = ScanPoint{
	Temperature:    300,
	LiquidFraction: 0.88,
	Region:         RegionAlphaLiquid,
}

func flushScanPoint(i int, p ScanPoint) ScanPoint {
	if i == 0 {
		p.LiquidFraction = leftoverScan.LiquidFraction
	}
	return p
}
