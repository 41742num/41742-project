package main

func main() {
	// 调用功能函数
	arrayA := generateAndPrintArray()
	sliceA, sliceB := createSlicesFromArray(arrayA)
	modifyAndRangeSlices(sliceA, sliceB) //当修改切片内容的时候，底层的数组元素相应的也会发生改变

	appendToSliceA(arrayA, sliceA, sliceB) //切片长度=数组原始长度+append元素个数；切片容量=原容量*倍数（扩容策略
	appendToSliceB(arrayA, sliceA, sliceB) //同上，原切片容量9，先扩容9*2-18

	copySliceAToSliceB(sliceA, sliceB) //切片copy时，会覆盖底层数组的元素，所以sliceA本身在copy完也会改变
	copyArrayToSliceCAndModify(arrayA) //通过copy过来的slice并不共享底层数组；如果是直接通过切片表达式arrayA[:]得出的就会共享

	ProcessFrame(arrayA)

}
